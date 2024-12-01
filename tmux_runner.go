package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type TmuxRunner struct {
	session    string
	doneString string
	interval   time.Duration
	debugLog   bool
}

var _ Runner = (*TmuxRunner)(nil)

func newTmuxRunner(debugLog bool) *TmuxRunner {
	return &TmuxRunner{
		session:    "git-exec-session",
		doneString: "git-exec-done",
		debugLog:   debugLog,
	}
}

func (x *TmuxRunner) Run(c *Command) (rerr error) {
	if _, err := exec.LookPath("tmux"); err != nil {
		return fmt.Errorf("tmux is not installed. Please install tmux. See https://github.com/tmux/tmux/wiki/Installing")
	}

	ch := make(chan error)
	go func() {
		ch <- x.tmuxNewSession()
	}()

	inputs := append(append(c.Envs, c.Args...), "; echo "+x.doneString)
	if err := x.tmuxSendKeys(inputs...); err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp("", "git-exec-pipe-pane")
	if err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	for {
		found, err := x.findDoneStringFromPipePane(tmpFile.Name())
		if err != nil {
			return err
		}
		if found {
			break
		}
	}

	output, err := x.tmuxCapturePane()
	if err != nil {
		return err
	}
	c.Output = output

	if err := x.killSession(); err != nil {
		return err
	}

	if err := <-ch; err != nil {
		return err
	}

	return nil
}

func (x *TmuxRunner) tmux(subcommand string, args ...string) error {
	arguments := append([]string{subcommand}, args...)
	cmd := exec.Command("tmux", arguments...)
	return cmd.Run()
}

func (x *TmuxRunner) tmuxNewSession() error {
	return x.tmux("new", "-s", x.session)
}

func (x *TmuxRunner) tmuxSendKeys(args ...string) error {
	arguments := []string{
		"send-keys",
		"-t",
		x.session,
		strings.Join(args, " "),
		"C-m",
	}
	return x.tmux("send-keys", arguments...)
}

func (x *TmuxRunner) pipePane(args ...string) error {
	arguments := append([]string{
		"pipe-pane",
		"-t",
		x.session,
	}, args...)
	return x.tmux("pipe-pane", arguments...)
}

func (x *TmuxRunner) tmuxCapturePane() (string, error) {
	cmd := exec.Command("tmux", "capture-pane", "-t", x.session, "-pS", "-", "-e")
	b, err := cmd.Output()
	return string(b), err
}

func (x *TmuxRunner) killSession() error {
	return x.tmux("kill-session", "-t", x.session)
}

func singleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func (x *TmuxRunner) findDoneStringFromPipePane(tmpFileName string) (bool, error) {
	if err := x.pipePane("-o", singleQuote("cat > "+tmpFileName)); err != nil {
		return false, err
	}

	time.Sleep(x.interval)

	if err := x.pipePane(); err != nil {
		return false, err
	}

	tmpFile, err := os.Open(tmpFileName)
	if err != nil {
		return false, err
	}
	defer tmpFile.Close()

	b, err := io.ReadAll(tmpFile)
	if err != nil {
		return false, err
	}
	return bytes.Contains(b, []byte(x.doneString)), nil
}
