package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitToOptionsAndCommandArgs(t *testing.T) {
	patterns := []struct {
		args        []string
		options     []string
		commandArgs []string
	}{
		{
			[]string{"-a", "-b", "key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			[]string{"-a", "-b"},
			[]string{"key1=val1", "key2=val2", "command", "--arg1", "arg2"},
		},
		{
			[]string{"-a", "-b", "command", "--arg1", "arg2"},
			[]string{"-a", "-b"},
			[]string{"command", "--arg1", "arg2"},
		},
		{
			[]string{"-a", "-b", "key1=val1", "key2=val2", "command"},
			[]string{"-a", "-b"},
			[]string{"key1=val1", "key2=val2", "command"},
		},
		{
			[]string{"command", "--arg1", "arg2"},
			nil,
			[]string{"command", "--arg1", "arg2"},
		},
		{
			[]string{"--version"},
			[]string{"--version"},
			nil,
		},
		{
			[]string{"-h"},
			[]string{"-h"},
			nil,
		},
	}

	for i, ptn := range patterns {
		t.Run(fmt.Sprintf("pattern %d", i), func(t *testing.T) {
			options, commandArgs := splitToOptionsAndCommandArgs(ptn.args)
			assert.Equal(t, ptn.options, options)
			assert.Equal(t, ptn.commandArgs, commandArgs)
		})
	}

}
