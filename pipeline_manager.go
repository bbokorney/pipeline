package main

import "github.com/bbokorney/dockworker/client"

// Manager manages starting and running pipelines
type Manager interface {
	NotifyNewPipeline(pipeline Pipeline)
	Start()
	Stop()
}

// NewManager returns a new Manager
func NewManager(dwClient client.Client, updater Updater, webhookListener WebhookListener) Manager {
	return manager{
		dwClient:        dwClient,
		newPipelineChan: make(chan Pipeline, 100),
		updater:         updater,
		webhookListener: webhookListener,
	}
}

type manager struct {
	dwClient        client.Client
	newPipelineChan chan Pipeline
	updater         Updater
	webhookListener WebhookListener
}

func (m manager) NotifyNewPipeline(pipeline Pipeline) {
	m.newPipelineChan <- pipeline
}

func (m manager) Start() {
	// TODO: ensure only one backgroundWorker is running
	go m.backgroundWorker()
}

func (m manager) Stop() {
	// TODO: implement
}

func (m manager) backgroundWorker() {
	for {
		select {
		case p := <-m.newPipelineChan:
			go NewWorker(p, m.dwClient, m.webhookListener, m.updater).Run()
		}
	}
}
