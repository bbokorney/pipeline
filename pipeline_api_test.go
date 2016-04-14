package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	retryCount         = 20
	pipelineHostEnvKey = "PIPELINE_URL"
)

func contains(i int, slice []string) bool {
	for _, arg := range slice {
		if strconv.Itoa(i) == arg {
			return true
		}
	}
	return false
}

func TestLargeAPI(t *testing.T) {
	url := os.Getenv(pipelineHostEnvKey)
	if url == "" {
		t.Fatalf("Must specify %s", pipelineHostEnvKey)
	}
	pipelineURL := fmt.Sprintf("%s/%s", url, "pipelines")

	runlist := os.Args[1:]
	fmt.Println("Running tests", runlist)

	for i, tc := range apiTestCases {
		if len(runlist) > 0 && !contains(i, runlist) {
			continue
		}
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
		assert.Equal(t, expectedSteps[i].Status, actualSteps[i].Status, "Case %d, Step %d: Step statuses should match", tcNum, i)
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
		if step.Status != StatusNotRun {
			assert.Condition(t, func() bool { return step.StartTime.Before(step.EndTime) || step.StartTime.Equal(step.EndTime) },
				"Case %d, Step %d: Start time (%s) should be before or equal to end time (%s)", tcNum, i, step.StartTime, step.EndTime)
			assert.NotEqual(t, 0, step.StartTime.Unix(), "Case %d, Step %d: Start time (%s) should not be 0")
			assert.NotEqual(t, 0, step.EndTime.Unix(), "Case %d, Step %d: End time (%s) should not be 0")
			steps[step.Name] = *step
		}
	}

	for stepIndex, step := range actualSteps {
		if stepDone(*step) {
			for depIndex, dep := range step.After {
				assert.Condition(t, func() bool { return steps[dep].Status == StatusSuccessful },
					"Case %d, Step %d: Dep %d: End time of dependency should be before start time of step", tcNum, stepIndex, depIndex)
				assert.Condition(t, func() bool { return steps[dep].EndTime.Before(step.StartTime) },
					"Case %d, Step %d: Dep %d: End time of dependency (%s) should be before start time of step (%s)",
					tcNum, stepIndex, depIndex, steps[dep].EndTime, step.StartTime)
			}
		}
	}
}

func waitUntilDone(t *testing.T, tcNum int, pipelineURL string, pipelineID PipelineID) {
	for i := 0; i < retryCount; i++ {
		p := getPipeline(t, tcNum, pipelineURL, pipelineID)
		if p.Status != StatusRunning && p.Status != StatusQueued && p.Status != StatusStopping {
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
