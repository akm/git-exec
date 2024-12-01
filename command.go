package main

import (
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
