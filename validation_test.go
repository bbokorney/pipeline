package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	for i, tc := range validationTestCases {
		err := ValidatePipeline(tc.pipeline)
		assert.Equal(t, tc.err, err, "Case %d: Error should match", i)
	}
}

type validationTestCase struct {
	pipeline Pipeline
	err      error
}

var validationTestCases = []validationTestCase{
	validationTestCase{
		err: ErrMissingPipelineName,
		pipeline: Pipeline{
			Name: "",
			Steps: []Step{
				Step{
					Name:      "Test Step",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
			},
		},
	},
	validationTestCase{
		err: ErrNoSteps,
		pipeline: Pipeline{
			Name:  "Test Pipeline",
			Steps: []Step{},
		},
	},
	validationTestCase{
		err: ErrMissingStepName,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []Step{
				Step{
					Name:      "Test Step",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				Step{
					Name:      "",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
			},
		},
	},
	validationTestCase{
		err: ErrNonUniqueStepNames,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []Step{
				Step{
					Name:      "Test Step",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				Step{
					Name:      "Test Step",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
			},
		},
	},
	validationTestCase{
		err: ErrMissingImageName,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []Step{
				Step{
					Name:      "Test Step 1",
					ImageName: "",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
			},
		},
	},
	validationTestCase{
		err: ErrMissingCommands,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []Step{
				Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", ""},
				},
				Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
			},
		},
	},
	validationTestCase{
		err: ErrMissingCommands,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []Step{
				Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{},
				},
			},
		},
	},
	validationTestCase{
		err: ErrNonExistentStepDependency,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []Step{
				Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"not a step"},
				},
				Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
				},
			},
		},
	},
	validationTestCase{
		err: ErrCircularStepDependency,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []Step{
				Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 2"},
				},
				Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
					After:     []string{"Test Step 1"},
				},
			},
		},
	},
	validationTestCase{
		err: ErrCircularStepDependency,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []Step{
				Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 1"},
				},
				Step{
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
			Steps: []Step{
				Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 2"},
				},
				Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
				},
			},
		},
	},
	validationTestCase{
		err: ErrCircularStepDependency,
		pipeline: Pipeline{
			Name: "Test Pipeline",
			Steps: []Step{
				Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 2"},
				},
				Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
					After:     []string{"Test Step 3"},
				},
				Step{
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
			Steps: []Step{
				Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 2"},
				},
				Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
				},
				Step{
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
			Steps: []Step{
				Step{
					Name:      "Test Step 1",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
				},
				Step{
					Name:      "Test Step 2",
					ImageName: "someimage:1234",
					Cmds:      []Cmd{"cmd"},
				},
				Step{
					Name:      "Test Step 3",
					ImageName: "someimage:123",
					Cmds:      []Cmd{"cmd1", "cmd2"},
					After:     []string{"Test Step 1"},
				},
			},
		},
	},
}
