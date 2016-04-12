package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bbokorney/dockworker"
	"github.com/bbokorney/dockworker/client"
	"github.com/emicklei/go-restful"
	"github.com/kelseyhightower/envconfig"
	"github.com/pborman/uuid"
)

// Config represents the program's config
type Config struct {
	DockworkerURL string `default:"http://dockworker:4321"`
	BindAddress   string `default:"0.0.0.0"`
	BindPort      int    `default:"4322"`
	WebhookURL    string `default:"http://pipeline:4322/webhook"`
}

var config Config

func doInit() *restful.Container {
	if err := envconfig.Process("pipeline", &config); err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}
	// TODO: logging level from config
	log.SetLevel(log.DebugLevel)
	wsContainer := restful.NewContainer()
	wsContainer.Filter(globalLogging)
	dwClient := client.NewClient(config.DockworkerURL)
	webhookChan := make(chan dockworker.Job)
	webhookListener := NewWebhookListener(webhookChan, config.WebhookURL)
	webhookListener.Start()
	pipelineStore := NewPipelineStore()
	updater := NewUpdater(pipelineStore)
	manager := NewManager(dwClient, updater, webhookListener)
	manager.Start()
	pipelineService := NewPipelineService(pipelineStore, manager)
	pipelineAPI := NewPipelineAPI(pipelineService)
	webhookAPI := NewWebhookAPI(webhookChan)
	pipelineAPI.Register(wsContainer)
	webhookAPI.Register(wsContainer)
	return wsContainer
}

func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	reqID := uuid.New()
	log.Infof("%s %s %s", req.Request.Method, req.Request.URL, reqID)
	chain.ProcessFilter(req, resp)
	log.Infof("%d %s %s", resp.StatusCode(), req.Request.URL, reqID)
}
