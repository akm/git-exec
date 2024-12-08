package core

import (
	"os"
	"testing"

	"github.com/akm/git-exec/command"

	"github.com/stretchr/testify/assert"
)

func TestCommitMessage(t *testing.T) {
	type pattern struct {
		name     string
		expected string
		location string
		emoji    string
		prompt   string
		template string
		envs     []string
		args     []string
		output   string
	}

	type option = func(*pattern)
	location := func(v string) option { return func(p *pattern) { p.location = v } }
	emoji := func(v string) option { return func(p *pattern) { p.emoji = v } }
	prompt := func(v string) option { return func(p *pattern) { p.prompt = v } }
	template := func(v string) option { return func(p *pattern) { p.template = v } }
	envs := func(v []string) option { return func(p *pattern) { p.envs = v } }
	args := func(v []string) option { return func(p *pattern) { p.args = v } }
	output := func(v string) option { return func(p *pattern) { p.output = v } }

	newPattern := func(name, expected string, opts ...option) pattern {
		p := pattern{
			name:     name,
			expected: expected,
			location: "root",
		}
		for _, opt := range opts {
			opt(&p)
		}
		return p
	}

	patterns := []pattern{
		newPattern("the most simple pattern", "🤖 [root] $ cmd1", args([]string{"cmd1"})),
		newPattern("with location", "🤖 [root/sub1] $ cmd1", location("root/sub1"), args([]string{"cmd1"})),
		newPattern("with envs", "🤖 [root] $ key1=val1 key2=val2 key3=val3 cmd1", args([]string{"cmd1"}),
			envs([]string{"key1=val1", "key2=val2", "key3=val3"}),
		),
		newPattern("with output", "🤖 [root] $ cmd1\n\noutput1", args([]string{"cmd1"}), output("output1")),
		newPattern("with multiple lines output", "🤖 [root] $ cmd1\n\noutput1\noutput2", args([]string{"cmd1"}), output("output1\noutput2")),
		newPattern("with arguments", "🤖 [root] $ cmd1 foo bar", args([]string{"cmd1", "foo", "bar"})),
		newPattern("with emoji", "🏭 [root] $ cmd1", emoji("🏭"), args([]string{"cmd1"})),
		newPattern("with prompt", "🤖 [root] % cmd1", prompt("%"), args([]string{"cmd1"})),
		newPattern("with template", "🤖 $ cmd1 [root]", args([]string{"cmd1"}),
			template("{{.Emoji}} {{.Prompt}} {{.Command}} [{{.Location}}]"),
		),
		newPattern("with template and multiple lines output", "🤖 $ cmd1 [root]\n\noutput1\noutput2", args([]string{"cmd1"}),
			output("output1\noutput2"),
			template("{{.Emoji}} {{.Prompt}} {{.Command}} [{{.Location}}]"),
		),
	}

	setEnvTemporarily := func(key, value string) func() {
		var setBackup func()
		backup := os.Getenv(key)
		if backup == "" {
			setBackup = func() { os.Unsetenv(key) }
		} else {
			setBackup = func() { os.Setenv(key, backup) }
		}
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
		return setBackup
	}

	for _, ptn := range patterns {
		t.Run(ptn.name, func(t *testing.T) {
			defer setEnvTemporarily("GIT_EXEC_EMOJI", ptn.emoji)()
			defer setEnvTemporarily("GIT_EXEC_PROMPT", ptn.prompt)()
			defer setEnvTemporarily("GIT_EXEC_TEMPLATE", ptn.template)()

			command := &command.Command{Envs: ptn.envs, Args: ptn.args, Output: ptn.output}
			commitMsg := newCommitMessage(command, newOptions())
			commitMsg.Location = ptn.location

			actual, err := commitMsg.Build()
			assert.NoError(t, err)
			assert.Equal(t, ptn.expected, actual)
		})
	}
}
