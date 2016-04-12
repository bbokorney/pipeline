package main

import (
	"sync"

	"github.com/bbokorney/dockworker"
)

// WebhookListener handles sending webhook notifications
// to the appropriate listeners
type WebhookListener interface {
	Start()
	Stop()
	Register(chan dockworker.Job)
	Unregister(chan dockworker.Job)
	WebhookURL() string
}

// NewWebhookListener returns a new WebhookListener
func NewWebhookListener(webhookChan chan dockworker.Job, webhookURL string) WebhookListener {
	return &webhookListener{
		webhookChan: webhookChan,
		listeners:   make(map[chan dockworker.Job]bool),
		lock:        &sync.RWMutex{},
		webhookURL:  webhookURL,
	}
}

type webhookListener struct {
	webhookChan chan dockworker.Job
	listeners   map[chan dockworker.Job]bool
	lock        *sync.RWMutex
	webhookURL  string
}

func (wl *webhookListener) Start() {
	go wl.backgroundWorker()
}

func (wl *webhookListener) Stop() {
	// TODO: implement
}

func (wl *webhookListener) Register(listener chan dockworker.Job) {
	wl.lock.Lock()
	defer wl.lock.Unlock()
	wl.listeners[listener] = true
}

func (wl *webhookListener) Unregister(listener chan dockworker.Job) {
	wl.lock.Lock()
	defer wl.lock.Unlock()
	delete(wl.listeners, listener)
}

func (wl *webhookListener) WebhookURL() string {
	return wl.webhookURL
}

func (wl *webhookListener) backgroundWorker() {
	for job := range wl.webhookChan {
		wl.lock.RLock()
		for listener := range wl.listeners {
			go sendMessage(listener, job)
		}
		wl.lock.RUnlock()
	}
}

func sendMessage(dest chan dockworker.Job, job dockworker.Job) {
	dest <- job
}
