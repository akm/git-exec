package main

import (
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Envs   []string
	Args   []string
	Output string
}

func newCommand(args []string) *Command {
	envs, commandArgs := splitArgsToEnvsAndCommand(args)
	return &Command{
		Envs: envs,
		Args: commandArgs,
	}
}

func (c *Command) Run() error {
	cmd := exec.Command(c.Args[0], c.Args[1:]...)
	cmd.Env = append(os.Environ(), c.Envs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	c.Output = string(output)
	return nil
}

func splitArgsToEnvsAndCommand(args []string) ([]string, []string) {
	equalNotFound := false
	return splitStringsInto2(args, func(arg string) bool {
		if !equalNotFound && strings.Contains(arg, "=") {
			return true
		} else {
			equalNotFound = true
			return false
		}
	})
}
