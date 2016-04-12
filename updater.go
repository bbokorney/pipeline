package main

import log "github.com/Sirupsen/logrus"

// Updater handles updating pipelines and steps
type Updater interface {
	UpdatePipeline(pipeline Pipeline) error
}

// NewUpdater returns a new Updater
func NewUpdater(pipelineStore PipelineStore) Updater {
	return updater{
		pipelineStore: pipelineStore,
	}
}

type updater struct {
	pipelineStore PipelineStore
}

func (u updater) UpdatePipeline(p Pipeline) error {
	if err := u.pipelineStore.Update(p); err != nil {
		log.Errorf("Error updating pipeline status %d: %s", p.ID, err)
		return err
	}
	return nil
}
