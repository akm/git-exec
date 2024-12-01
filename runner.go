package main

type Runner interface {
	Run(c *Command) error
}
