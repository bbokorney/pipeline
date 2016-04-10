package main

// Pipeline is a set of Steps
type Pipeline struct {
	ID     PipelineID `json:"id"`
	Name   string     `json:"name"`
	Steps  []Step     `json:"steps"`
	Status Status     `json:"status"`
}

// Step is a step of a Pipeline
type Step struct {
	Name      string   `json:"name"`
	ImageName string   `json:"image"`
	Cmds      []Cmd    `json:"cmds"`
	After     []string `json:"after"`
	Status    Status   `json:"status"`
}

// PipelineID is and identifier for a Pipeline
type PipelineID int

// Cmd represents a command
type Cmd string

// Status represents the state of a
// Pipeline or Step
type Status string

const (
	// StatusQueued state indicates the job is queued waiting to be run
	StatusQueued Status = "queued"
	// StatusRunning state indicates the job is running
	StatusRunning Status = "running"
	// StatusSuccessful state indicates the job has completed successfully
	StatusSuccessful Status = "successful"
	// StatusFailed state indicates the job has completed with a failure
	StatusFailed Status = "failed"
	// StatusError state indicates the job could not be run properly
	StatusError Status = "error"
)
