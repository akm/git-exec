package command

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitArgsToEnvsAndCommand(t *testing.T) {
	patterns := []struct {
		args        []string
		envs        []string
		commandArgs []string
	}{
		{
			[]string{"key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			[]string{"key1=val1", "key2=val2"},
			[]string{"command", "--arg1", "arg2"},
		},
		{
			[]string{"key1=val1", "command", "--arg1", "arg2", "key2=val2"},
			[]string{"key1=val1"},
			[]string{"command", "--arg1", "arg2", "key2=val2"},
		},
		{
			[]string{"command", "--arg1", "arg2"},
			nil,
			[]string{"command", "--arg1", "arg2"},
		},
	}

	for i, ptn := range patterns {
		t.Run(fmt.Sprintf("pattern %d", i), func(t *testing.T) {
			envs, commandArgs := splitArgsToEnvsAndCommand(ptn.args)
			assert.Equal(t, ptn.envs, envs)
			assert.Equal(t, ptn.commandArgs, commandArgs)
		})
	}

}
