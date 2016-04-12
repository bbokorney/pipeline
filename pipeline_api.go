package main

import (
	"net/http"
	"strconv"

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
			response.WriteHeaderAndEntity(http.StatusNotFound, errorResponse("No pipeline with that ID"))
			return
		default:
			response.WriteHeaderAndEntity(http.StatusInternalServerError, errorResponse(err.Error()))
			return
		}
	}
	response.WriteHeaderAndEntity(http.StatusOK, pipeline)
}

func (api PipelineAPI) createPipeline(request *restful.Request, response *restful.Response) {
	pipeline := &Pipeline{}
	err := request.ReadEntity(pipeline)
	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	p, err := api.pipelineService.Add(*pipeline)
	if err != nil {
		if isValidationError(err) {
			response.WriteHeaderAndEntity(http.StatusBadRequest, errorResponse(err.Error()))
			return
		}
		response.WriteHeaderAndEntity(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}
	response.WriteHeaderAndEntity(http.StatusCreated, p)
}
