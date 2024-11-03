package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	defaultExpected := func() *Options {
		return &Options{
			Emoji:    "🤖",
			Prompt:   "$",
			Template: `{{.Emoji}} [{{.Location}}] {{.Prompt}} {{.Command}}`,
		}
	}
	version := func(v bool) func(*Options) { return func(o *Options) { o.Version = v } }
	help := func(v bool) func(*Options) { return func(o *Options) { o.Help = v } }
	directory := func(v string) func(*Options) { return func(o *Options) { o.Directory = v } }
	emoji := func(v string) func(*Options) { return func(o *Options) { o.Emoji = v } }
	prompt := func(v string) func(*Options) { return func(o *Options) { o.Prompt = v } }
	template := func(v string) func(*Options) { return func(o *Options) { o.Template = v } }
	newExpected := func(opts ...func(*Options)) *Options {
		o := defaultExpected()
		for _, opt := range opts {
			opt(o)
		}
		return o
	}

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
			newExpected(),
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"--emoji", "🏭", "command"},
			newExpected(emoji("🏭")),
			[]string{"command"},
			"",
		},
		{
			map[string]string{"GIT_EXEC_EMOJI": "🏭"},
			[]string{"command"},
			newExpected(emoji("🏭")),
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"--prompt", "%", "command"},
			newExpected(prompt("%")),
			[]string{"command"},
			"",
		},
		{
			map[string]string{"GIT_EXEC_PROMPT": "%"},
			[]string{"command"},
			newExpected(prompt("%")),
			[]string{"command"},
			"",
		},
		{
			// Both --prompt and GIT_EXEC_PROMPT are given, --prompt should take precedence
			map[string]string{"GIT_EXEC_PROMPT": "%"},
			[]string{"--prompt", ">", "command"},
			newExpected(prompt(">")),
			[]string{"command"},
			"",
		},
		{
			map[string]string{"GIT_EXEC_TEMPLATE": "{{.Emoji}} {{.Prompt}} {{.Command}} [{{.Location}}]"},
			[]string{"command"},
			newExpected(template("{{.Emoji}} {{.Prompt}} {{.Command}} [{{.Location}}]")),
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"--template", "{{.Emoji}} {{.Prompt}} {{.Command}} [{{.Location}}]", "command"},
			newExpected(template("{{.Emoji}} {{.Prompt}} {{.Command}} [{{.Location}}]")),
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"-h", "-v", "key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			newExpected(help(true), version(true)),
			[]string{"key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			"",
		},
		{
			map[string]string{},
			[]string{"-C", "foo", "key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			newExpected(directory("foo")),
			[]string{"key1=val1", "key2=val2", "command", "--arg1", "arg2"},
			"",
		},
		{
			map[string]string{},
			[]string{"-v", "-h", "command", "--arg1", "arg2"},
			newExpected(help(true), version(true)),
			[]string{"command", "--arg1", "arg2"},
			"",
		},
		{
			map[string]string{},
			[]string{"--directory", "bar", "command", "--arg1", "arg2"},
			newExpected(directory("bar")),
			[]string{"command", "--arg1", "arg2"},
			"",
		},
		{
			map[string]string{},
			[]string{"-h", "-v", "key1=val1", "key2=val2", "command"},
			newExpected(help(true), version(true)),
			[]string{"key1=val1", "key2=val2", "command"},
			"",
		},
		{
			map[string]string{},
			[]string{"command", "--arg1", "arg2"},
			newExpected(),
			[]string{"command", "--arg1", "arg2"},
			"",
		},
		{
			map[string]string{},
			[]string{"--directory", "baz", "command"},
			newExpected(directory("baz")),
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"--directory", "baz", "-v", "command"},
			newExpected(directory("baz"), version(true)),
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"--directory", "--help", "-v", "command"},
			newExpected(directory("--help"), version(true)),
			[]string{"command"},
			"",
		},
		{
			map[string]string{},
			[]string{"--version"},
			newExpected(version(true)),
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
			newExpected(help(true)),
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
