package main

// PipelineService manages Pipelines
type PipelineService interface {
	Add(pipeline Pipeline) (Pipeline, error)
	Find(ID PipelineID) (Pipeline, error)
}

// NewPipelineService returns a new PipelineService
func NewPipelineService(pipelineStore PipelineStore) PipelineService {
	return pipelineService{
		pipelineStore: pipelineStore,
	}
}

type pipelineService struct {
	pipelineStore PipelineStore
}

// Add creates a new Pipeline
func (service pipelineService) Add(pipeline Pipeline) (Pipeline, error) {
	// add all of the steps
	pipeline.Status = StatusQueued
	p, err := service.pipelineStore.Add(pipeline)
	if err != nil {
		return Pipeline{}, err
	}
	return p, nil
}

func (service pipelineService) Find(ID PipelineID) (Pipeline, error) {
	return service.pipelineStore.Find(ID)
}
