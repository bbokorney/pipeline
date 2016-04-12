package main

import (
	"testing"

	"github.com/bbokorney/dockworker"
	"github.com/stretchr/testify/assert"
)

func TestConvertCmds(t *testing.T) {
	cmds := []Cmd{"touch file.txt", "run.sh start stop", "ls"}
	expected := []dockworker.Cmd{
		dockworker.Cmd{"touch", "file.txt"},
		dockworker.Cmd{"run.sh", "start", "stop"},
		dockworker.Cmd{"ls"},
	}
	converted := convertCmds(cmds)
	assert.Equal(t, expected, converted, "Converted commands should match")
}
