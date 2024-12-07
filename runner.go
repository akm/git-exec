package main

import (
	"github.com/akm/git-exec/command"
)

type Runner interface {
	Run(c *command.Command) error
}
