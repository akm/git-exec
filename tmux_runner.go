package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
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

	time.Sleep(1 * time.Second)

	defer func() {
		if err := x.killSession(); err != nil && rerr == nil {
			rerr = err
		}
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

	if err := <-ch; err != nil {
		return err
	}

	return nil
}

func (x *TmuxRunner) tmux(subcommand string, args ...string) error {
	arguments := append([]string{subcommand}, args...)
	cmd := exec.Command("tmux", arguments...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	slog.Debug("tmux", "args", arguments)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tmux %s: %w", strings.Join(arguments, " "), err)
	}
	return nil
}

func (x *TmuxRunner) tmuxNewSession() error {
	return x.tmux("new", "-s", x.session)
}

func (x *TmuxRunner) tmuxSendKeys(args ...string) error {
	arguments := []string{
		"-t",
		x.session,
		singleQuote(strings.Join(args, " ")),
		"C-m",
	}
	return x.tmux("send-keys", arguments...)
}

func (x *TmuxRunner) pipePane(args ...string) error {
	arguments := append([]string{
		"-t",
		x.session,
	}, args...)
	return x.tmux("pipe-pane", arguments...)
}

func (x *TmuxRunner) tmuxCapturePane() (string, error) {
	cmd := exec.Command("tmux", "capture-pane", "-t", x.session, "-pS", "-", "-e")
	slog.Debug("tmux capture-pane", "args", cmd.Args)
	b, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("tmux capture-pane: %w", err)
	}
	return string(b), nil
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

func init() {
	f, err := os.Open("debug.log")
	if err == nil {
		panic("failed to open debug.log")
	}
	logger := slog.New(slog.NewTextHandler(f, nil))
	slog.SetDefault(logger)
}
