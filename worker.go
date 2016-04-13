package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bbokorney/dockworker"
	"github.com/bbokorney/dockworker/client"
)

// Worker runs a Pipeline
type Worker interface {
	Run()
}

// NewWorker returns a new worker
func NewWorker(pipeline Pipeline, dwClient client.Client, webhookListener WebhookListener, updater Updater) Worker {
	webhookChan := make(chan dockworker.Job)
	webhookListener.Register(webhookChan)
	steps := make(map[string]*Step)
	for _, step := range pipeline.Steps {
		steps[step.Name] = step
	}
	return &worker{
		pipeline:        &pipeline,
		dwClient:        dwClient,
		webhookListener: webhookListener,
		updater:         updater,
		webhookChan:     webhookChan,
		steps:           steps,
		runningJobs:     make(map[dockworker.JobID]int),
	}
}

type worker struct {
	pipeline        *Pipeline
	dwClient        client.Client
	webhookListener WebhookListener
	updater         Updater
	webhookChan     chan dockworker.Job
	steps           map[string]*Step
	runningJobs     map[dockworker.JobID]int
}

func (w *worker) Run() {
	defer w.cleanup()
	log.Infof("Starting run of pipeline %d", w.pipeline.ID)
	w.updatePipelineStatus(StatusRunning)
	// initialize ourselves with a step to run
	if err := w.doRun(); err != nil {
		w.updatePipelineStatus(StatusError)
		log.Errorf("Failed to run pipeline %d: %s", w.pipeline.ID, err)
		return
	}
	log.Infof("Successful run of pipeline %d", w.pipeline.ID)
}

func (w *worker) doRun() error {
	// start the steps with no dependencies
	if err := w.runReadySteps(); err != nil {
		return err
	}
	for {
		select {
		case jobUpdate := <-w.webhookChan:
			log.Debugf("Received job update %+v", jobUpdate)
			// we've received a job update
			done, err := w.handleUpdate(jobUpdate)
			if err != nil {
				return err
			}
			if done {
				return nil
			}
		}
	}
}

func (w *worker) handleUpdate(job dockworker.Job) (done bool, err error) {
	// check if this one of our jobs
	if _, contains := w.runningJobs[job.ID]; !contains {
		// this isn't a job for this pipeline
		return false, nil
	}

	stepIndex := w.runningJobs[job.ID]
	delete(w.runningJobs, job.ID)
	w.pipeline.Steps[stepIndex].StartTime = job.StartTime
	w.pipeline.Steps[stepIndex].EndTime = job.EndTime
	w.pipeline.Steps[stepIndex].JobURL = fmt.Sprintf("%s/jobs/%d", w.dwClient.BaseURL(), job.ID)

	log.Debugf("Job %d has status %s", job.ID, job.Status)
	// set the status of the step
	switch job.Status {
	case dockworker.JobStatusFailed:
		w.pipeline.Steps[stepIndex].Status = StatusFailed
	case dockworker.JobStatusError:
		w.pipeline.Steps[stepIndex].Status = StatusError
	case dockworker.JobStatusStopped:
		w.pipeline.Steps[stepIndex].Status = StatusStopped
	case dockworker.JobStatusSuccessful:
		w.pipeline.Steps[stepIndex].Status = StatusSuccessful
	}

	if w.pipeline.Status != StatusStopping && w.pipeline.Steps[stepIndex].Status != StatusSuccessful {
		// this is the first detection of failure
		// We need to start cleaning up
		w.pipeline.Status = StatusStopping
		log.Debugf("Pipeline %d has status %s", w.pipeline.ID, StatusStopping)
		w.stopRunningJobs()
		w.setQueuedToNotRun()
	}

	// check if we're stopping the pipeline
	if w.pipeline.Status == StatusStopping {
		done = false
		if len(w.runningJobs) == 0 {
			// no jobs left running, we can exit
			w.pipeline.Status = StatusFailed
			log.Debugf("Pipeline %d has status %s and 0 running jobs", w.pipeline.ID, StatusFailed)
			done = true
		}
		w.saveUpdatedPipeline()
		return done, nil
	}

	// now we know the job was successful and
	// the pipeline is still running
	w.saveUpdatedPipeline()

	// this job finishing could have satisfied the
	// dependencies of another waiting step
	// start any such steps
	if err := w.runReadySteps(); err != nil {
		return true, err
	}

	// check if all the steps are done
	if w.allStepsDone() {
		w.pipeline.Status = StatusSuccessful
		w.saveUpdatedPipeline()
		return true, nil
	}
	return false, nil
}

func (w *worker) stopRunningJobs() {
	for jobID, stepIndex := range w.runningJobs {
		log.Debugf("Stopping job %d for step %d", jobID, stepIndex)
		if err := w.dwClient.StopJob(dockworker.JobID(jobID)); err != nil {
			log.Errorf("Error stopping job %d: %s", jobID, err)
		}
		w.pipeline.Steps[stepIndex].Status = StatusStopped
	}
}

func (w *worker) setQueuedToNotRun() {
	for _, step := range w.pipeline.Steps {
		if step.Status == StatusQueued {
			step.Status = StatusNotRun
		}
	}
}

func (w *worker) runReadySteps() error {
	log.Debug("Running ready steps")
	for i, step := range w.pipeline.Steps {
		if stepDoneOrRunning(*step) {
			continue
		}
		if w.dependenciesDone(*step) {
			log.Debugf("Running step %+v", step)
			if err := w.runStep(step, i); err != nil {
				return err
			}
			log.Debugf("Done starting step %+v", step)
		}
	}
	return nil
}

func (w *worker) runStep(step *Step, stepIndex int) error {
	job := dockworker.Job{
		ImageName:  step.ImageName,
		Cmds:       convertCmds(step.Cmds),
		WebhookURL: w.webhookListener.WebhookURL(),
	}
	createdJob, err := w.dwClient.CreateJob(job)
	if err != nil {
		log.Errorf("Failed to create job %+v for pipeline %d: %s", job, w.pipeline.ID, err)
		return err
	}
	log.Debugf("Job started %+v", createdJob)
	w.runningJobs[createdJob.ID] = stepIndex
	step.Status = StatusRunning
	w.saveUpdatedPipeline()
	return nil
}

func (w *worker) updatePipelineStatus(status Status) {
	w.pipeline.Status = status
	w.saveUpdatedPipeline()
}

func (w *worker) saveUpdatedPipeline() {
	if err := w.updater.UpdatePipeline(*w.pipeline); err != nil {
		log.Errorf("Failed to update status of pipeline %d", w.pipeline.ID)
	}
}

func (w *worker) dependenciesDone(step Step) bool {
	for _, dep := range step.After {
		if w.steps[dep].Status != StatusSuccessful {
			return false
		}
	}
	return true
}

func (w *worker) allStepsDone() bool {
	for _, step := range w.pipeline.Steps {
		if !stepDone(*step) {
			return false
		}
	}
	return true
}

func stepDone(step Step) bool {
	return step.Status == StatusSuccessful ||
		step.Status == StatusFailed ||
		step.Status == StatusError
}

func stepRunning(step Step) bool {
	return step.Status == StatusRunning
}

func stepDoneOrRunning(step Step) bool {
	return stepDone(step) || stepRunning(step)
}

func (w *worker) cleanup() {
	// unregister and empty the webhook channel
	go func() {
		for _ = range w.webhookChan {
		}
	}()
	w.webhookListener.Unregister(w.webhookChan)
}

func convertCmds(cmds []Cmd) []dockworker.Cmd {
	var converted []dockworker.Cmd
	for _, c := range cmds {
		converted = append(converted, strings.Split(string(c), " "))
	}
	return converted
}
