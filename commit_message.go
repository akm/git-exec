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

func newCommitMessage(location, env, command string) *commitMessage {
	return &commitMessage{
		Env:      env,
		Emoji:    getEnvString("GIT_EXEC_EMOJI", "🤖"),
		Location: location,
		Prompt:   getEnvString("GIT_EXEC_PROMPT", "$"),
		Command:  command,
	}
}

const defaultHeadTemplateSurce = `{{.Emoji}} @{{.Location}} {{.Prompt}} {{.Command}}`

func newTemplate() (*template.Template, error) {
	source := getEnvString("GIT_EXEC_TEMPLATE", defaultHeadTemplateSurce)
	return template.New("commitMessage").Parse(source + "\n\n{{.Body}}\n")
}

func buildCommitMessage(command *Command) (string, error) {
	commandParts := make([]string, len(command.Args))
	for i, arg := range command.Args {
		if strings.Contains(arg, " ") && !(strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'")) {
			commandParts[i] = fmt.Sprintf("'%s'", arg)
		} else {
			commandParts[i] = arg
		}
	}

	location, err := getLocation()
	if err != nil {
		return "", err
	}

	envs := []string{}
	if len(command.Envs) > 0 {
		envs = append(envs, strings.Join(command.Envs, " "))
	}

	msg := newCommitMessage(location, strings.Join(envs, " "), strings.Join(commandParts, " "))
	msg.Body = command.Output

	tmpl, err := newTemplate()
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, msg); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getLocation() (string, error) {
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
