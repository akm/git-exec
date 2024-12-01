package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Envs   []string
	Args   []string
	Output string

	debugLog bool
}

func newCommand(args []string) *Command {
	envs, commandArgs := splitArgsToEnvsAndCommand(args)
	return &Command{
		Envs: envs,
		Args: commandArgs,
	}
}

func (c *Command) EnableDebugLog() {
	c.debugLog = true
}

func (c *Command) Run() (rerr error) {
	cmd := exec.Command(c.Args[0], c.Args[1:]...)
	cmd.Env = append(os.Environ(), c.Envs...)
	cmd.Stdin = os.Stdin
	var buf bytes.Buffer

	stdoutWriter, stdoutTd, err := c.newOutputWriter(os.Stdout, &buf, "./stdout.log")
	if err != nil {
		return nil
	}
	defer func() {
		if err := stdoutTd(); err != nil && rerr == nil {
			rerr = err
		}
	}()

	stderrWriter, stderrTd, err := c.newOutputWriter(os.Stderr, &buf, "./stderr.log")
	if err != nil {
		return nil
	}
	defer func() {
		if err := stderrTd(); err != nil && rerr == nil {
			rerr = err
		}
	}()

	cmd.Stdout = stdoutWriter
	cmd.Stderr = stderrWriter
	if err := cmd.Run(); err != nil {
		return err
	}
	c.Output = buf.String()
	return nil
}

func (c *Command) newOutputWriter(original io.Writer, buf *bytes.Buffer, debugLogFile string) (io.Writer, func() error, error) {
	if !c.debugLog {
		return io.MultiWriter(os.Stdout, buf), func() error { return nil }, nil
	}
	logFile, err := os.OpenFile(debugLogFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}
	return io.MultiWriter(original, buf, logFile), logFile.Close, nil
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
