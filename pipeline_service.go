package main

// PipelineService manages Pipelines
type PipelineService interface {
	Add(pipeline Pipeline) (Pipeline, error)
	Find(ID int) (Pipeline, error)
}

// NewPipelineService returns a new PipelineService
func NewPipelineService(pipelineStore PipelineStore, stepStore StepStore) PipelineService {
	return pipelineService{
		pipelineStore: pipelineStore,
		stepStore:     stepStore,
	}
}

type pipelineService struct {
	pipelineStore PipelineStore
	stepStore     StepStore
}

// Add creates a new Pipeline
func (service pipelineService) Add(pipeline Pipeline) (Pipeline, error) {
	// add all of the steps
	var steps []Step
	for _, step := range pipeline.Steps {
		s, err := service.stepStore.Add(step)
		if err != nil {
			return Pipeline{}, err
		}
		steps = append(steps, s)
	}
	pipeline.Steps = steps
	p, err := service.pipelineStore.Add(pipeline)
	if err != nil {
		return Pipeline{}, err
	}
	return p, nil
}

func (service pipelineService) Find(ID int) (Pipeline, error) {
	return service.pipelineStore.Find(ID)
}
