package main

// Pipeline is a set of Steps
type Pipeline struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
}

// Step is a step of a Pipeline
type Step struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	ImageName string  `json:"image"`
	Cmds      []Cmd   `json:"cmds"`
	Inputs    []Input `json:"inputs"`
}

// InputType is an input type
type InputType string

const (
	// PrevStep is an input from a previous step
	PrevStep InputType = "prev_step"
	// Repo is an input from a repo
	Repo InputType = "repo"
)

// Input is an input to a step
type Input map[string]string
