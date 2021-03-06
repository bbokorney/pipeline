package main

import "time"

// Pipeline is a set of Steps
type Pipeline struct {
	ID     PipelineID `json:"id"`
	Name   string     `json:"name"`
	Steps  []*Step    `json:"steps"`
	Status Status     `json:"status"`
}

// TODO: investigate omit if empty struct tags

// Step is a step of a Pipeline
type Step struct {
	Name      string            `json:"name"`
	ImageName string            `json:"image"`
	Cmds      []Cmd             `json:"cmds"`
	Env       map[string]string `json:"env"`
	After     []string          `json:"after"`
	JobURL    string            `json:"job_url"`
	Status    Status            `json:"status"`
	StartTime time.Time         `json:"start_time"`
	EndTime   time.Time         `json:"end_time"`
}

// PipelineID is and identifier for a Pipeline
type PipelineID int

// Cmd represents a command
type Cmd string

// Status represents the state of a
// Pipeline or Step
type Status string

const (
	// StatusQueued state indicates the job is queued waiting to be run
	StatusQueued Status = "queued"
	// StatusRunning state indicates the job is running
	StatusRunning Status = "running"
	// StatusSuccessful state indicates the job has completed successfully
	StatusSuccessful Status = "successful"
	// StatusFailed state indicates the job has completed with a failure
	StatusFailed Status = "failed"
	// StatusStopping state indicates the job is stopping
	StatusStopping Status = "stopping"
	// StatusError state indicates the job could not be run properly
	StatusError Status = "error"
	// StatusNotRun state indicates the job was not run
	StatusNotRun Status = "not-run"
	// StatusStopped state indicates the job was stoped
	StatusStopped Status = "stopped"
)

// NotRunTime represents the time for steps which have not been started or ended
var NotRunTime = time.Unix(0, 0)
