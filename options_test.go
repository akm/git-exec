package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	patterns := []struct {
		envs        map[string]string
		args        []string
		options     *Options
		commandArgs []string
		error       string
	}{
		{
			map[string]string{},
			[]string{"command"},
			&Options{},
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"--emoji", "🏭", "command"},
			&Options{Emoji: "🏭"},
			[]string{"command"},
			"",
		},
		{
			map[string]string{"GIT_EXEC_EMOJI": "🏭"},
			[]string{"command"},
			&Options{Emoji: "🏭"},
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"--prompt", "%", "command"},
			&Options{Prompt: "%"},
			[]string{"command"},
			"",
		},
		{
			map[string]string{"GIT_EXEC_PROMPT": "%"},
			[]string{"command"},
			&Options{Prompt: "%"},
			[]string{"command"},
			"",
		},
		{
			// Both --prompt and GIT_EXEC_PROMPT are given, --prompt should take precedence
			map[string]string{"GIT_EXEC_PROMPT": "%"},
			[]string{"--prompt", ">", "command"},
			&Options{Prompt: ">"},
			[]string{"command"},
			"",
		},
		{
			map[string]string{"GIT_EXEC_TEMPLATE": "{{.Emoji}} {{.Prompt}} {{.Command}} [{{.Location}}]"},
			[]string{"command"},
			&Options{Template: "{{.Emoji}} {{.Prompt}} {{.Command}} [{{.Location}}]"},
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"--template", "{{.Emoji}} {{.Prompt}} {{.Command}} [{{.Location}}]", "command"},
			&Options{Template: "{{.Emoji}} {{.Prompt}} {{.Command}} [{{.Location}}]"},
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"-h", "-v", "key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			&Options{Help: true, Version: true},
			[]string{"key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			"",
		},
		{
			map[string]string{},
			[]string{"-C", "foo", "key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			&Options{Directory: "foo"},
			[]string{"key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			"",
		},
		{
			map[string]string{},
			[]string{"-v", "-h", "command", "--arg1", "arg2"},
			&Options{Version: true, Help: true},
			[]string{"command", "--arg1", "arg2"},
			"",
		},
		{
			map[string]string{},
			[]string{"--directory", "bar", "command", "--arg1", "arg2"},
			&Options{Directory: "bar"},
			[]string{"command", "--arg1", "arg2"},
			"",
		},
		{
			map[string]string{},
			[]string{"-h", "-v", "key1=val1", "key2=val2", "command"},
			&Options{Help: true, Version: true},
			[]string{"key1=val1", "key2=val2", "command"},
			"",
		},
		{
			map[string]string{},
			[]string{"command", "--arg1", "arg2"},
			&Options{},
			[]string{"command", "--arg1", "arg2"},
			"",
		},
		{
			map[string]string{},
			[]string{"--directory", "baz", "command"},
			&Options{Directory: "baz"},
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"--directory", "baz", "-v", "command"},
			&Options{Directory: "baz", Version: true},
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"--directory", "--help", "-v", "command"},
			&Options{Directory: "--help", Version: true},
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"--version"},
			&Options{Version: true},
			[]string{},
			"",
		},
		{
			map[string]string{},
			[]string{"--directory"},
			nil,
			nil,
			"no value given for option --directory",
		},
		{
			map[string]string{},
			[]string{"-h"},
			&Options{Help: true},
			[]string{},
			"",
		},
	}

	for _, ptn := range patterns {
		t.Run(fmt.Sprintf("pattern %+v %v", ptn.envs, ptn.args), func(t *testing.T) {
			for key, value := range ptn.envs {
				key := key
				envBackup := os.Getenv(key)
				os.Setenv(key, value)
				defer func() { os.Setenv(key, envBackup) }()
			}

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
