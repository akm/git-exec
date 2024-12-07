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
	equalNotFound := false
	var a, b []string
	for _, arg := range args {
		if !equalNotFound && strings.Contains(arg, "=") {
			a = append(a, arg)
		} else {
			b = append(b, arg)
		}
	}
	return a, b
}
