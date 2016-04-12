package main

import "time"

// PipelineService manages Pipelines
type PipelineService interface {
	Add(pipeline Pipeline) (Pipeline, error)
	Find(ID PipelineID) (Pipeline, error)
}

// NewPipelineService returns a new PipelineService
func NewPipelineService(pipelineStore PipelineStore, manager Manager) PipelineService {
	return pipelineService{
		pipelineStore: pipelineStore,
		manager:       manager,
	}
}

type pipelineService struct {
	pipelineStore PipelineStore
	manager       Manager
}

// Add creates a new Pipeline
func (service pipelineService) Add(pipeline Pipeline) (Pipeline, error) {
	if err := ValidatePipeline(pipeline); err != nil {
		return Pipeline{}, err
	}

	pipeline.Status = StatusQueued
	for _, step := range pipeline.Steps {
		step.Status = StatusQueued
		step.StartTime = time.Unix(0, 0)
		step.EndTime = time.Unix(0, 0)
	}
	p, err := service.pipelineStore.Add(pipeline)
	if err != nil {
		return Pipeline{}, err
	}
	service.manager.NotifyNewPipeline(p)
	return p, nil
}

func (service pipelineService) Find(ID PipelineID) (Pipeline, error) {
	return service.pipelineStore.Find(ID)
}
