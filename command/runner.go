package command

type Runner interface {
	Run(c *Command) error
}
