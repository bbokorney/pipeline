package main

import "sync"

type JobStore interface {
	Add(job Job) (Job, error)
	Find(ID int) (Job, error)
}

func NewJobStore() JobStore {
	return &inMemJobStore{
		lock:   &sync.RWMutex{},
		nextID: 0,
		data:   make(map[int]Job),
	}
}

type inMemJobStore struct {
	lock   *sync.RWMutex
	nextID int
	data   map[int]Job
}

func (store *inMemJobStore) Add(j Job) (Job, error) {
	store.lock.Lock()
	defer store.lock.Unlock()
	id := store.nextID
	store.data[id] = j
	store.nextID = store.nextID + 1
	return Job{
		ID:         id,
		StepID:     j.StepID,
		BranchName: j.BranchName,
		CommitHash: j.CommitHash,
		step:       j.step,
	}, nil
}

func (store inMemJobStore) Find(ID int) (Job, error) {
	store.lock.RLock()
	defer store.lock.RUnlock()
	s, ok := store.data[ID]
	if !ok {
		return Job{}, ErrNotFound
	}
	return s, nil
}
