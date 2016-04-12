package main

import (
	"errors"
	"sync"
)

// PipelineStore stores pipelines
type PipelineStore interface {
	Add(pipeline Pipeline) (Pipeline, error)
	Find(ID PipelineID) (Pipeline, error)
	Update(p Pipeline) error
}

var (
	// ErrNotFound indicates an item not found
	ErrNotFound = errors.New("Item not found")
)

// NewPipelineStore returns a new PipelineStore
func NewPipelineStore() PipelineStore {
	return &inMemPipelineStore{
		lock:   &sync.RWMutex{},
		nextID: 0,
		data:   make(map[PipelineID]Pipeline),
	}
}

type inMemPipelineStore struct {
	lock   *sync.RWMutex
	nextID PipelineID
	data   map[PipelineID]Pipeline
}

func (store *inMemPipelineStore) Add(p Pipeline) (Pipeline, error) {
	store.lock.Lock()
	defer store.lock.Unlock()
	p.ID = store.nextID
	store.data[p.ID] = p
	store.nextID = store.nextID + 1
	return p, nil
}

func (store inMemPipelineStore) Find(ID PipelineID) (Pipeline, error) {
	store.lock.RLock()
	defer store.lock.RUnlock()
	p, ok := store.data[ID]
	if !ok {
		return Pipeline{}, ErrNotFound
	}
	return p, nil
}

func (store *inMemPipelineStore) Update(p Pipeline) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	store.data[p.ID] = p
	return nil
}
