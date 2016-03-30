package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	wsContainer := initWSContainer()

	ts := httptest.NewServer(wsContainer)
	defer ts.Close()

	pipelineURL := fmt.Sprintf("%s/%s", ts.URL, "pipelines")

	resp, err := http.Post(pipelineURL, "application/json", strings.NewReader(examplePipelineBody))
	if err != nil {
		t.Errorf("Error sending post request: %s", err)
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Status code should be 201")

	pipelinePOST := decodeBody(t, resp.Body)

	resp, err = http.Get(fmt.Sprintf("%s/%d", pipelineURL, pipelinePOST.ID))
	if err != nil {
		t.Errorf("Error sending get request: %s", err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code should be 200")
	pipelineGET := decodeBody(t, resp.Body)

	assert.Equal(t, pipelineGET, pipelinePOST)
}

func decodeBody(t *testing.T, respBody io.ReadCloser) *Pipeline {
	body, err := ioutil.ReadAll(respBody)
	defer respBody.Close()
	if err != nil {
		t.Errorf("Error reading response body: %s", err)
	}

	pipeline := &Pipeline{}
	err = json.Unmarshal(body, pipeline)
	if err != nil {
		t.Errorf("Error decoding response body: %s", err)
	}

	return pipeline
}

const examplePipelineBody = `{
  "name": "Pipeline Name",
  "steps": [
    {
      "name": "Step Name",
      "image": "ubuntu:14.04",
      "cmds": [
        ["ls", "-la"],
        ["touch", "hello.txt"],
        ["ls", "-la"]
      ],
      "input": {
        "repo": {
          "dir": "/path/to/repo",
          "branch": "branch-name"
        },
        "prev_steps": [
          {
            "step_name": "step name",
            "dir": "/some/path"
          },
          {
            "step_name": "another step name",
            "dir": "/some/other/path"
          }
        ]
      },
      "output": {
        "dir": "/output/path"
      }
    }
  ]
}
`
