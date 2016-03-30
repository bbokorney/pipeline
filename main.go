package main

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

func main() {
	wsContainer := initWSContainer()
	log.Fatal(http.ListenAndServe(":4321", wsContainer))
}

func initWSContainer() *restful.Container {
	pipelineAPI := initPipelineAPI()
	wsContainer := restful.NewContainer()
	pipelineAPI.Register(wsContainer)
	return wsContainer
}

func initPipelineAPI() PipelineAPI {
	pipelineStore := NewPipelineStore()
	stepStore := NewStepStore()
	pipelineService := NewPipelineService(pipelineStore, stepStore)
	return NewPipelineAPI(pipelineService)
}
