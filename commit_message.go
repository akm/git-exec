package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type commitMessage struct {
	Env      string
	Emoji    string
	Location string
	Prompt   string
	Command  string
	Body     string
}

func newCommitMessage(command *Command) *commitMessage {
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
		Env:     strings.Join(envs, " "),
		Emoji:   getEnvString("GIT_EXEC_EMOJI", "🤖"),
		Prompt:  getEnvString("GIT_EXEC_PROMPT", "$"),
		Command: strings.Join(commandParts, " "),
		Body:    command.Output,
	}
}

const defaultHeadTemplateSurce = `{{.Emoji}} @{{.Location}} {{.Prompt}} {{.Command}}`

func (*commitMessage) newTemplate() (*template.Template, error) {
	source := getEnvString("GIT_EXEC_TEMPLATE", defaultHeadTemplateSurce)
	return template.New("commitMessage").Parse(source + "\n\n{{.Body}}\n")
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
