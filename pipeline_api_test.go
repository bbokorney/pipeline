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

	for i, tc := range apiTestCases {
		resp, err := http.Post(pipelineURL, "application/json", strings.NewReader(tc.exampleBody))
		if err != nil {
			t.Errorf("Case %d: Error sending post request: %s", i, err)
		}
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Case %d: Status code should be 201", i)

		pipelinePOST := decodeBody(t, i, resp.Body)

		resp, err = http.Get(fmt.Sprintf("%s/%d", pipelineURL, pipelinePOST.ID))
		if err != nil {
			t.Errorf("Case %d: Error sending get request: %s", i, err)
		}
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Case %d: Status code should be 200", i)
		pipelineGET := decodeBody(t, i, resp.Body)

		assert.Equal(t, pipelineGET, pipelinePOST)
	}
}

func decodeBody(t *testing.T, tcNum int, respBody io.ReadCloser) *Pipeline {
	body, err := ioutil.ReadAll(respBody)
	defer respBody.Close()
	if err != nil {
		t.Errorf("Case %d: Error reading response body: %s", tcNum, err)
	}

	pipeline := &Pipeline{}
	err = json.Unmarshal(body, pipeline)
	if err != nil {
		t.Errorf("Error decoding response body: %s", err)
	}

	return pipeline
}

type apiTestCase struct {
	exampleBody string
}

var apiTestCases = []apiTestCase{
	apiTestCase{
		exampleBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "Step Name",
	      "image": "ubuntu:14.04",
	      "cmds": [
	        "ls -la",
	        "touch hello.txt",
	        "ls -la"
	      ]
			}
	  ]
	}`,
	},
}
