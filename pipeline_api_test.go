package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	retryCount         = 20
	pipelineHostEnvKey = "PIPELINE_URL"
)

func TestAPI(t *testing.T) {
	url := os.Getenv(pipelineHostEnvKey)
	if url == "" {
		t.Fatalf("Must specify %s", pipelineHostEnvKey)
	}
	pipelineURL := fmt.Sprintf("%s/%s", url, "pipelines")

	for i, tc := range apiTestCases {
		resp, err := http.Post(pipelineURL, "application/json", strings.NewReader(tc.requestBody))
		if err != nil {
			t.Errorf("Case %d: Error sending post request: %s", i, err)
		}
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Case %d: Status code should be 201", i)

		pipelinePOST := decodeBody(t, i, resp.Body)
		assert.Equal(t, tc.pipeline.Name, pipelinePOST.Name, "Case %d: Pipeline name should be unchanged", i)
		compareStepData(t, i, tc.pipeline.Steps, pipelinePOST.Steps)
		assert.Equal(t, StatusQueued, pipelinePOST.Status, "Case %d: Status should be queued", i)

		waitUntilDone(t, i, pipelineURL, pipelinePOST.ID)
		pipelineGET := getPipeline(t, i, pipelineURL, pipelinePOST.ID)
		assert.Equal(t, tc.pipeline.Status, pipelineGET.Status, "Case %d: Status should match", i)
		compareStepStatuses(t, i, tc.pipeline.Steps, pipelineGET.Steps)
		compareStepJobURLs(t, i, tc.pipeline.Steps, pipelineGET.Steps)
		validateStepTimestampsAndDependencies(t, i, pipelineGET.Steps)
	}
}

func compareStepData(t *testing.T, tcNum int, expectedSteps []*Step, actualSteps []*Step) {
	assert.Equal(t, len(expectedSteps), len(actualSteps), "Case %d: Length of steps should match", tcNum)
	for i := range expectedSteps {
		assert.Equal(t, expectedSteps[i].Name, actualSteps[i].Name, "Case %d, Step %d: Step name should be unchanged", tcNum, i)
		assert.Equal(t, expectedSteps[i].ImageName, actualSteps[i].ImageName, "Case %d, Step %d: Step image should be unchanged", tcNum, i)
		assert.Equal(t, expectedSteps[i].After, actualSteps[i].After, "Case %d, Step %d: Step dependencies should be unchanged", tcNum, i)
		assert.Equal(t, expectedSteps[i].Cmds, actualSteps[i].Cmds, "Case %d, Step %d: Step cmds should be unchanged", tcNum, i)
	}
}

func compareStepStatuses(t *testing.T, tcNum int, expectedSteps []*Step, actualSteps []*Step) {
	assert.Equal(t, len(expectedSteps), len(actualSteps), "Case %d: Length of steps should match", tcNum)
	for i := range expectedSteps {
		assert.Equal(t, expectedSteps[i].Status, actualSteps[i].Status, "Case %d, Step %d: Step name should be unchanged", tcNum, i)
	}
}

func compareStepJobURLs(t *testing.T, tcNum int, expectedSteps []*Step, actualSteps []*Step) {
	assert.Equal(t, len(expectedSteps), len(actualSteps), "Case %d: Length of steps should match", tcNum)
	for i := range expectedSteps {
		if expectedSteps[i].Status != StatusQueued && expectedSteps[i].Status != StatusNotRun {
			assert.NotEmpty(t, actualSteps[i].JobURL, "Case %d, Step %d: JobURL should be set", tcNum, i)
		}
	}
}

func validateStepTimestampsAndDependencies(t *testing.T, tcNum int, actualSteps []*Step) {
	steps := make(map[string]Step)
	for i, step := range actualSteps {
		assert.Condition(t, func() bool { return step.StartTime.Before(step.EndTime) },
			"Case %d, Step %d: Start time should be before end time", tcNum, i)
		steps[step.Name] = *step
	}

	for stepIndex, step := range actualSteps {
		if stepDone(*step) {
			for depIndex, dep := range step.After {
				assert.Condition(t, func() bool { return steps[dep].Status == StatusSuccessful },
					"Case %d, Step %d: Dep %d: End time of dependency should be before start time of step", tcNum, stepIndex, depIndex)
				assert.Condition(t, func() bool { return steps[dep].EndTime.Before(step.StartTime) },
					"Case %d, Step %d: Dep %d: End time of dependency should be before start time of step", tcNum, stepIndex, depIndex)
			}
		}
	}
}

func waitUntilDone(t *testing.T, tcNum int, pipelineURL string, pipelineID PipelineID) {
	for i := 0; i < retryCount; i++ {
		p := getPipeline(t, tcNum, pipelineURL, pipelineID)
		if p.Status != "running" && p.Status != "queued" {
			return
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatalf("Case %d: Waitied too long for pipeline to complete", tcNum)
}

func getPipeline(t *testing.T, tcNum int, pipelineURL string, ID PipelineID) *Pipeline {
	resp, err := http.Get(fmt.Sprintf("%s/%d", pipelineURL, ID))
	if err != nil {
		t.Errorf("Case %d: Error sending get request: %s", tcNum, err)
	}
	return decodeBody(t, tcNum, resp.Body)
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

func getLogs(t *testing.T, tcNum int, url string) string {
	resp, err := http.Get(fmt.Sprintf("%s/logs", url))
	if err != nil {
		t.Errorf("Case %d: Error getting logs: %s", tcNum, err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Case %d: Expected code getting logs %d but received %d", tcNum, http.StatusCreated, resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Case %d: Error reading log response body %s", tcNum, err)
	}
	return string(respBody)
}

type apiTestCase struct {
	requestBody string
	pipeline    Pipeline
}

var apiTestCases = []apiTestCase{
	apiTestCase{
		requestBody: `{
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
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "Step Name",
					ImageName: "ubuntu:14.04",
					Cmds: []Cmd{
						"ls -la",
						"touch hello.txt",
						"ls -la",
					},
					Status: StatusSuccessful,
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
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
			},
			{
	      "name": "Other Step Name",
	      "image": "ubuntu:14.04",
	      "cmds": [
	        "ls -la",
	        "touch hello.txt",
	        "ls -la"
	      ]
			}
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "Step Name",
					ImageName: "ubuntu:14.04",
					Cmds: []Cmd{
						"ls -la",
						"touch hello.txt",
						"ls -la",
					},
					Status: StatusSuccessful,
				},
				&Step{
					Name:      "Other Step Name",
					ImageName: "ubuntu:14.04",
					Cmds: []Cmd{
						"ls -la",
						"touch hello.txt",
						"ls -la",
					},
					Status: StatusSuccessful,
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step1"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1"]
	    },
			{
	      "name": "step3",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step2"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step1"},
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step2"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1"]
	    },
			{
	      "name": "step3",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step1"},
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step1"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step3",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1", "step2"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step1", "step2"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1"]
	    },
			{
	      "name": "step3",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step1"},
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
			},
		},
	},
	// Failure tests
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "Step Name",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls notafile.txt"]
			}
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusFailed,
			Steps: []*Step{
				&Step{
					Name:      "Step Name",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls notafile.txt"},
					Status:    StatusFailed,
				},
			},
		},
	},
}
