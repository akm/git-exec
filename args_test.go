package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitToOptionsAndCommandArgs(t *testing.T) {
	patterns := []struct {
		args        []string
		options     Options
		commandArgs []string
		error       string
	}{
		{
			[]string{"-h", "-v", "key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			Options{&Option{Type: optHelp}, &Option{Type: optVersion}},
			[]string{"key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			"",
		},
		{
			[]string{"-v", "-h", "command", "--arg1", "arg2"},
			Options{&Option{Type: optVersion}, &Option{Type: optHelp}},
			[]string{"command", "--arg1", "arg2"},
			"",
		},
		{
			[]string{"-h", "-v", "key1=val1", "key2=val2", "command"},
			Options{&Option{Type: optHelp}, &Option{Type: optVersion}},
			[]string{"key1=val1", "key2=val2", "command"},
			"",
		},
		{
			[]string{"command", "--arg1", "arg2"},
			Options{},
			[]string{"command", "--arg1", "arg2"},
			"",
		},
		{
			[]string{"--version"},
			Options{&Option{Type: optVersion}},
			[]string{},
			"",
		},
		{
			[]string{"-h"},
			Options{&Option{Type: optHelp}},
			[]string{},
			"",
		},
	}

	for i, ptn := range patterns {
		t.Run(fmt.Sprintf("pattern %d", i), func(t *testing.T) {
			options, commandArgs, err := splitToOptionsAndCommandArgs(ptn.args)
			assert.Equal(t, ptn.options, options)
			assert.Equal(t, ptn.commandArgs, commandArgs)
			if ptn.error == "" {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, ptn.error, err.Error())
			}
		})
	}

}
