package main

import (
	"errors"
	"sync"
)

// PipelineStore stores pipelines
type PipelineStore interface {
	Add(pipeline Pipeline) (Pipeline, error)
	Find(ID int) (Pipeline, error)
}

type StepStore interface {
	Add(step Step) (Step, error)
	Find(ID int) (Step, error)
}

var (
	ErrNotFound = errors.New("Item not found")
)

func NewPipelineStore() PipelineStore {
	return &inMemPipelineStore{
		lock:   &sync.RWMutex{},
		nextID: 0,
		data:   make(map[int]Pipeline),
	}
}

func NewStepStore() StepStore {
	return &inMemStepStore{
		lock:   &sync.RWMutex{},
		nextID: 0,
		data:   make(map[int]Step),
	}
}

type inMemPipelineStore struct {
	lock   *sync.RWMutex
	nextID int
	data   map[int]Pipeline
}

func (store *inMemPipelineStore) Add(p Pipeline) (Pipeline, error) {
	store.lock.Lock()
	defer store.lock.Unlock()
	id := store.nextID
	store.data[id] = p
	store.nextID = store.nextID + 1
	return Pipeline{
		ID:    id,
		Name:  p.Name,
		Steps: p.Steps,
	}, nil
}

func (store inMemPipelineStore) Find(ID int) (Pipeline, error) {
	store.lock.RLock()
	defer store.lock.RUnlock()
	p, ok := store.data[ID]
	if !ok {
		return Pipeline{}, ErrNotFound
	}
	return p, nil
}

type inMemStepStore struct {
	lock   *sync.RWMutex
	nextID int
	data   map[int]Step
}

func (store *inMemStepStore) Add(s Step) (Step, error) {
	store.lock.Lock()
	defer store.lock.Unlock()
	id := store.nextID
	store.data[id] = s
	store.nextID = store.nextID + 1
	return Step{
		ID:        id,
		Name:      s.Name,
		ImageName: s.ImageName,
		Cmds:      s.Cmds,
		Inputs:    s.Inputs,
	}, nil
}

func (store inMemStepStore) Find(ID int) (Step, error) {
	store.lock.RLock()
	defer store.lock.RUnlock()
	s, ok := store.data[ID]
	if !ok {
		return Step{}, ErrNotFound
	}
	return s, nil
}
