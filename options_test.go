package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	patterns := []struct {
		args        []string
		options     *Options
		commandArgs []string
		error       string
	}{
		{
			[]string{"-h", "-v", "key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			&Options{Help: true, Version: true},
			[]string{"key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			"",
		},
		{
			[]string{"-C", "foo", "key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			&Options{Directory: "foo"},
			[]string{"key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			"",
		},
		{
			[]string{"-v", "-h", "command", "--arg1", "arg2"},
			&Options{Version: true, Help: true},
			[]string{"command", "--arg1", "arg2"},
			"",
		},
		{
			[]string{"--directory", "bar", "command", "--arg1", "arg2"},
			&Options{Directory: "bar"},
			[]string{"command", "--arg1", "arg2"},
			"",
		},
		{
			[]string{"-h", "-v", "key1=val1", "key2=val2", "command"},
			&Options{Help: true, Version: true},
			[]string{"key1=val1", "key2=val2", "command"},
			"",
		},
		{
			[]string{"command", "--arg1", "arg2"},
			&Options{},
			[]string{"command", "--arg1", "arg2"},
			"",
		},
		{
			[]string{"--directory", "baz", "command"},
			&Options{Directory: "baz"},
			[]string{"command"},
			"",
		},
		{
			[]string{"--directory", "baz", "-v", "command"},
			&Options{Directory: "baz", Version: true},
			[]string{"command"},
			"",
		},
		{
			[]string{"--directory", "--help", "-v", "command"},
			&Options{Directory: "--help", Version: true},
			[]string{"command"},
			"",
		},
		{
			[]string{"--version"},
			&Options{Version: true},
			[]string{},
			"",
		},
		{
			[]string{"--directory"},
			nil,
			nil,
			"no value given for option --directory",
		},
		{
			[]string{"-h"},
			&Options{Help: true},
			[]string{},
			"",
		},
	}

	for i, ptn := range patterns {
		t.Run(fmt.Sprintf("pattern %d", i), func(t *testing.T) {
			options, commandArgs, err := parseOptions(ptn.args)
			assert.Equal(t, ptn.options, options)
			assert.Equal(t, ptn.commandArgs, commandArgs)
			if ptn.error == "" {
				assert.Nil(t, err)
			} else {
				if assert.NotNil(t, err) {
					assert.Equal(t, ptn.error, err.Error())
				}
			}
		})
	}

}
