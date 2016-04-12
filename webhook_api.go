package main

import (
	"net/http"

	"github.com/bbokorney/dockworker"
	"github.com/emicklei/go-restful"
)

// WebhookAPI is the webhook receiver API
type WebhookAPI struct {
	webhookChan chan dockworker.Job
}

// NewWebhookAPI returns a new WebhookAPI
func NewWebhookAPI(webhookChan chan dockworker.Job) WebhookAPI {
	return WebhookAPI{
		webhookChan: webhookChan,
	}
}

// Register adds the routes to the web service container
func (api WebhookAPI) Register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/webhook").
		Consumes(restful.MIME_JSON)

	ws.Route(ws.POST("").To(api.handleWebhook).
		Operation("handleWebhook").
		Reads(dockworker.Job{}))

	container.Add(ws)
}

func (api WebhookAPI) handleWebhook(request *restful.Request, response *restful.Response) {
	job := &dockworker.Job{}
	err := request.ReadEntity(job)
	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	go func() {
		api.webhookChan <- *job
	}()

	response.WriteHeader(http.StatusAccepted)
}
