package main

import "fmt"

var (
	// ErrMissingPipelineName indicates a pipeline name is missing
	ErrMissingPipelineName = fmt.Errorf("Must specify a pipeline name")
	// ErrMissingStepName indicates a step name is missing
	ErrMissingStepName = fmt.Errorf("Must specify a step name")
	// ErrNoSteps indicates a step name is missing
	ErrNoSteps = fmt.Errorf("Must specify at least one step")
	// ErrNonUniqueStepNames indicates not all step names are unique
	ErrNonUniqueStepNames = fmt.Errorf("All step names must be unique")
	// ErrMissingImageName indicates an image name is missing
	ErrMissingImageName = fmt.Errorf("Must specify an image name")
	// ErrMissingCommands indicates a step name is missing
	ErrMissingCommands = fmt.Errorf("Must specify a command or list of commands")
	// ErrNonExistentStepDependency indicates a step name is missing
	ErrNonExistentStepDependency = fmt.Errorf("All step dependencies must exist")
	// ErrCircularStepDependency indicates a step name is missing
	ErrCircularStepDependency = fmt.Errorf("Must have no circular dependencies between steps")
)

// ValidationError represents a pipeline validation error
type ValidationError struct {
	err error
}

func (ve ValidationError) Error() string {
	return ve.err.Error()
}

// ValidatePipeline checks that a pipeline is valid
// checks all the necessary fields aren't blank
// checks that either cmd or cmds only are specified
// checks the dependency graph of steps
// * all references are to steps which are specified
// * no circular references
func ValidatePipeline(pipeline Pipeline) error {
	return runValidations(pipeline)
}

func runValidations(pipeline Pipeline) error {
	for _, v := range validations {
		if err := v(pipeline); err != nil {
			return err
		}
	}
	return nil
}

type validation func(pipeline Pipeline) error

var validations = []validation{
	func(pipeline Pipeline) error {
		if pipeline.Name == "" {
			return ErrMissingPipelineName
		}
		return nil
	},
	func(pipeline Pipeline) error {
		if len(pipeline.Steps) < 1 {
			return ErrNoSteps
		}
		return nil
	},
	func(pipeline Pipeline) error {
		for _, step := range pipeline.Steps {
			if step.Name == "" {
				return ErrMissingStepName
			}
		}
		return nil
	},
	func(pipeline Pipeline) error {
		for _, step := range pipeline.Steps {
			if step.ImageName == "" {
				return ErrMissingImageName
			}
		}
		return nil
	},
	func(pipeline Pipeline) error {
		for _, step := range pipeline.Steps {
			if len(step.Cmds) < 1 {
				return ErrMissingCommands
			}
			// return error if any Cmds were specified
			// as blank
			for _, cmd := range step.Cmds {
				if cmd == "" {
					return ErrMissingCommands
				}
			}
		}
		return nil
	},
	func(pipeline Pipeline) error {
		steps := make(map[string]bool)
		for _, step := range pipeline.Steps {
			if _, contains := steps[step.Name]; contains {
				return ErrNonUniqueStepNames
			}
			steps[step.Name] = true
		}
		return nil
	},
	func(pipeline Pipeline) error {
		steps := make(map[string]Step)
		for _, step := range pipeline.Steps {
			steps[step.Name] = *step
		}
		// do a BFS through the graph
		overallVisited := make(map[string]bool)
		// continue until we've visted all steps
		for len(overallVisited) < len(steps) {
			visited := make(map[string]bool)
			q := queue{}
			// start the search with a node we haven't visited yet
			for step := range steps {
				if _, contains := overallVisited[step]; !contains {
					q.push(step)
					break
				}
			}
			for q.len() > 0 {
				curr := q.pop()
				// check if we've visited this step during this traversal
				if _, contains := visited[curr]; contains {
					return ErrCircularStepDependency
				}
				// check that this step is legit
				if _, contains := steps[curr]; !contains {
					return ErrNonExistentStepDependency
				}
				// mark this step as visited
				overallVisited[curr] = true
				visited[curr] = true
				for _, dep := range steps[curr].After {
					q.push(dep)
				}
			}
		}

		return nil
	},
}

type queue struct {
	slice []string
}

func (q *queue) len() int {
	return len(q.slice)
}

func (q *queue) push(item string) {
	q.slice = append(q.slice, item)
}

func (q *queue) pop() string {
	if len(q.slice) < 1 {
		return ""
	}
	head := q.slice[0]
	q.slice = q.slice[1:len(q.slice)]
	return head
}
