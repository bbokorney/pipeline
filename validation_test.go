package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmallValidation(t *testing.T) {
	for i, tc := range validationTestCases {
		err := ValidatePipeline(tc.pipeline)
		assert.Equal(t, tc.err, err, "Case %d: Error should match (error %s)", i, err)
	}
}

type validationTestCase struct {
	pipeline Pipeline
	err      error
}

var validationTestCases = []validationTestCase{
	validationTestCase{
		err: ValidationError{ErrMissingPipelineName},
		pipeline: Pipeline{
			Name: "",
			Steps: []*Step{
				&Step{
					Name:      "Test Step",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
			},
		},
	},
	validationTestCase{
		err: ValidationError{ErrNoSteps},
		pipeline: Pipeline{
			Name:  "Test Pipeline",
			Steps: []*Step{},
		},
	},
	validationTestCase{
		err: ValidationError{ErrMissingStepName},
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "Test Step",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				&Step{
					Name:      "",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
			},
		},
	},
	validationTestCase{
		err: ValidationError{ErrNonUniqueStepNames},
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "Test Step",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				&Step{
					Name:      "Test Step",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
			},
		},
	},
	validationTestCase{
		err: ValidationError{ErrMissingImageName},
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "Test Step 1",
					ImageName: "",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				&Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
			},
		},
	},
	validationTestCase{
		err: ValidationError{ErrMissingCommands},
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", ""},
				},
				&Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
			},
		},
	},
	validationTestCase{
		err: ValidationError{ErrMissingCommands},
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				&Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{},
				},
			},
		},
	},
	validationTestCase{
		err: ValidationError{ErrNonExistentStepDependency},
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"not a step"},
				},
				&Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
				},
			},
		},
	},
	validationTestCase{
		err: ValidationError{ErrCircularStepDependency},
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 2"},
				},
				&Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
					After:     []string{"Test Step 1"},
				},
			},
		},
	},
	validationTestCase{
		err: ValidationError{ErrCircularStepDependency},
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 1"},
				},
				&Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
				},
			},
		},
	},
	validationTestCase{
		err: nil,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 2"},
				},
				&Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
				},
			},
		},
	},
	validationTestCase{
		err: ValidationError{ErrCircularStepDependency},
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 2"},
				},
				&Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
					After:     []string{"Test Step 3"},
				},
				&Step{
					Name:      "Test Step 3",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 1"},
				},
			},
		},
	},
	validationTestCase{
		err: nil,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 2"},
				},
				&Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
				},
				&Step{
					Name:      "Test Step 3",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 1"},
				},
			},
		},
	},
	validationTestCase{
		err: nil,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				&Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
				},
				&Step{
					Name:      "Test Step 3",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 1"},
				},
			},
		},
	},
	validationTestCase{
		err: nil,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				&Step{
					Name:      "step2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
				},
				&Step{
					Name:      "step3",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"step1", "step2"},
				},
			},
		},
	},
	validationTestCase{
		err: nil,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				&Step{
					Name:      "step2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
					After:     []string{"step1"},
				},
				&Step{
					Name:      "step3",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"step1"},
				},
			},
		},
	},
	validationTestCase{
		err: nil,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				&Step{
					Name:      "step2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
					After:     []string{"step1"},
				},
				&Step{
					Name:      "step3",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"step1"},
				},
				&Step{
					Name:      "step4",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"step2", "step3"},
				},
			},
		},
	},
	validationTestCase{
		err: ValidationError{ErrCircularStepDependency},
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"step1"},
				},
			},
		},
	},
	validationTestCase{
		err: ValidationError{ErrCircularStepDependency},
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				&Step{
					Name:      "step2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
					After:     []string{"step1", "step5"},
				},
				&Step{
					Name:      "step3",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"step1"},
				},
				&Step{
					Name:      "step4",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"step2", "step3"},
				},
				&Step{
					Name:      "step5",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"step2", "step3"},
				},
			},
		},
	},
	validationTestCase{
		err: ValidationError{ErrCircularStepDependency},
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusFailed,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					After:     []string{"step2"},
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					After:     []string{"step4"},
				},
				&Step{
					Name:      "step4",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					After:     []string{"step5"},
				},
				&Step{
					Name:      "step5",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					After:     []string{"step3"},
				},
			},
		},
	},
}
