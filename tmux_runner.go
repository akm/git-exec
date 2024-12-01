package main

type TmuxRunner struct {
	debugLog bool
}

var _ Runner = (*TmuxRunner)(nil)

func (x *TmuxRunner) Run(c *Command) (rerr error) {
	panic("implement me")
}
