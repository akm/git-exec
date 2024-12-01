package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type StandardRunner struct {
	debugLog bool
}

var _ Runner = (*StandardRunner)(nil)

func newStandardRunner(debugLog bool) *StandardRunner {
	return &StandardRunner{
		debugLog: debugLog,
	}
}

func (x *StandardRunner) Run(c *Command) (rerr error) {
	cmd := exec.Command(c.Args[0], c.Args[1:]...)
	cmd.Env = append(os.Environ(), c.Envs...)
	cmd.Stdin = os.Stdin
	var buf bytes.Buffer

	stdoutWriter, stdoutTd, err := x.newOutputWriter(os.Stdout, &buf, "./stdout.log")
	if err != nil {
		return nil
	}
	defer func() {
		if err := stdoutTd(); err != nil && rerr == nil {
			rerr = err
		}
	}()

	stderrWriter, stderrTd, err := x.newOutputWriter(os.Stderr, &buf, "./stderr.log")
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

func (x *StandardRunner) newOutputWriter(original io.Writer, buf *bytes.Buffer, debugLogFile string) (io.Writer, func() error, error) {
	if !x.debugLog {
		return io.MultiWriter(os.Stdout, buf), func() error { return nil }, nil
	}
	logFile, err := os.OpenFile(debugLogFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}
	return io.MultiWriter(original, buf, logFile), logFile.Close, nil
}
