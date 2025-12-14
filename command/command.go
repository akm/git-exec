package command

import (
	"strings"
)

type Command struct {
	Envs   []string
	Args   []string
	Output string
}

func NewCommand(args []string) *Command {
	envs, commandArgs := splitArgsToEnvsAndCommand(args)
	return &Command{
		Envs: envs,
		Args: commandArgs,
	}
}

func splitArgsToEnvsAndCommand(args []string) ([]string, []string) {
	notEnvFound := false
	var a, b []string
	for _, arg := range args {
		if !notEnvFound && strings.Contains(arg, "=") {
			a = append(a, arg)
		} else {
			b = append(b, arg)
			notEnvFound = true
		}
	}
	return a, b
}
