package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	commandParts := make([]string, len(command.Args))
	for i, arg := range command.Args {
		if strings.Contains(arg, " ") && !(strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'")) {
			commandParts[i] = fmt.Sprintf("'%s'", arg)
		} else {
			commandParts[i] = arg
		}
	}

	envs := []string{}
	if len(command.Envs) > 0 {
		envs = append(envs, strings.Join(command.Envs, " "))
	}

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
	location, err := m.getLocation()
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

	return buf.String(), nil
}

func (*commitMessage) getLocation() (string, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	rootDir, err := gitRootDir()
	if err != nil {
		return "", err
	}
	rootDirName := filepath.Base(rootDir)

	relPath, err := filepath.Rel(rootDir, curDir)
	if err != nil {
		return "", err
	}

	if relPath == "." {
		return rootDirName, nil
	} else if strings.HasPrefix(relPath, "./") {
		return rootDirName + relPath[1:], nil
	} else if strings.HasPrefix(relPath, "/") {
		return relPath, nil
	} else {
		return rootDirName + "/" + relPath, nil
	}
}

func gitRootDir() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
