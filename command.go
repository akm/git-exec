package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Envs   []string
	Args   []string
	Output *bytes.Buffer
}

func newCommand(args []string) *Command {
	envs, commandArgs := splitArgsToEnvsAndCommand(args)
	return &Command{
		Envs:   envs,
		Args:   commandArgs,
		Output: &bytes.Buffer{},
	}
}

func (c *Command) Run() error {
	cmd := exec.Command(c.Args[0], c.Args[1:]...)
	cmd.Env = append(os.Environ(), c.Envs...)
	cmd.Stdout = c.Output
	cmd.Stderr = c.Output
	return cmd.Run()
}

func splitArgsToEnvsAndCommand(args []string) ([]string, []string) {
	var envs []string
	var command []string
	equalNotFound := false
	for _, arg := range args {
		if !equalNotFound && strings.Contains(arg, "=") {
			envs = append(envs, arg)
		} else {
			equalNotFound = true
			command = append(command, arg)
		}
	}
	return envs, command
}
