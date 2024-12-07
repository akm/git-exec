package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/akm/git-exec/command"
)

type commitMessage struct {
	Env      string
	Emoji    string
	Location string
	Prompt   string
	Template string
	Command  string
	Body     string
}

func newCommitMessage(command *command.Command, options *Options) *commitMessage {
	argParts := make([]string, len(command.Args))
	for i, arg := range command.Args {
		if strings.Contains(arg, " ") && !(strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'")) {
			argParts[i] = fmt.Sprintf("'%s'", arg)
		} else {
			argParts[i] = arg
		}
	}

	envs := []string{}
	if len(command.Envs) > 0 {
		envs = append(envs, strings.Join(command.Envs, " "))
	}

	commandParts := append(envs, argParts...)

	return &commitMessage{
		Env:      strings.Join(envs, " "),
		Emoji:    options.Emoji,
		Prompt:   options.Prompt,
		Template: options.Template,
		Command:  strings.Join(commandParts, " "),
		Body:     command.Output,
	}
}

func (m *commitMessage) newTemplate() (*template.Template, error) {
	return template.New("commitMessage").Parse(m.Template + "\n\n{{.Body}}\n")
}

func (m *commitMessage) Build() (string, error) {
	location, err := getLocation()
	if err != nil {
		return "", err
	}
	m.Location = location

	tmpl, err := m.newTemplate()
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, m); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}
