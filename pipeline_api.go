package main

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

// PipelineAPI is the Pipeline management API
type PipelineAPI struct {
	pipelineService PipelineService
}

// NewPipelineAPI returns a new PipelineAPI
func NewPipelineAPI(pipelineService PipelineService) PipelineAPI {
	return PipelineAPI{
		pipelineService: pipelineService,
	}
}

// Register adds the routes to the web service container
func (api PipelineAPI) Register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/pipelines").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{id}").To(api.findPipeline).
		Operation("findPipeline").
		Param(ws.PathParameter("id", "id of pipeline").DataType("int")).
		Writes(Pipeline{}))

	ws.Route(ws.POST("").To(api.createPipeline).
		Operation("createPipeline").
		Reads(Pipeline{}))

	container.Add(ws)
}

func (api PipelineAPI) findPipeline(request *restful.Request, response *restful.Response) {
	id, err := strconv.Atoi(request.PathParameter("id"))
	if err != nil {
		response.WriteHeaderAndEntity(http.StatusNotFound, errorResponse("ID must be int"))
		return
	}
	pipelineID := PipelineID(id)

	pipeline, err := api.pipelineService.Find(pipelineID)
	if err != nil {
		switch err {
		case ErrNotFound:
			logAndRespondError(response, http.StatusNotFound, err)
			return
		default:
			logAndRespondError(response, http.StatusInternalServerError, err)
			return
		}
	}
	response.WriteHeaderAndEntity(http.StatusOK, pipeline)
}

func (api PipelineAPI) createPipeline(request *restful.Request, response *restful.Response) {
	pipeline := &Pipeline{}
	err := request.ReadEntity(pipeline)
	if err != nil {
		logAndRespondError(response, http.StatusInternalServerError, err)
		return
	}

	p, err := api.pipelineService.Add(*pipeline)
	if err != nil {
		if isValidationError(err) {
			logAndRespondError(response, http.StatusBadRequest, err)
			return
		}
		logAndRespondError(response, http.StatusInternalServerError, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusCreated, p)
}

func logAndRespondError(response *restful.Response, status int, err error) {
	log.Infof("Error response %d %s", status, err)
	response.WriteHeaderAndEntity(status, errorResponse(err.Error()))
}
