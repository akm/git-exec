package main

import (
	"bytes"
	"os"
	"os/exec"
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
