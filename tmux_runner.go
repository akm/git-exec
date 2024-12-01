package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type TmuxRunner struct {
	session     string
	pipeLogFile string
	interval    time.Duration
	debugLog    bool

	completeString string
	errorString    string
}

var _ Runner = (*TmuxRunner)(nil)

func newTmuxRunner(debugLog bool) *TmuxRunner {
	doneStringPrefix := "git-exec-done"
	return &TmuxRunner{
		session:        "git-exec-session",
		pipeLogFile:    filepath.Join(os.TempDir(), "trace.log"),
		debugLog:       debugLog,
		interval:       1_000 * time.Millisecond,
		completeString: doneStringPrefix + "-complete",
		errorString:    doneStringPrefix + "-error",
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

	sessionKilled := false
	defer func() {
		if !sessionKilled {
			if err := x.killSession(); err != nil && rerr == nil {
				rerr = err
			}
		}
	}()

	inputs := append(append(c.Envs, c.Args...), " && echo "+x.completeString+" || echo "+x.errorString)
	if err := x.tmuxSendKeys(inputs...); err != nil {
		return err
	}

	defer func() {
		if err := os.Remove(x.pipeLogFile); err != nil && rerr == nil {
			rerr = err
		}
	}()

	if err := x.startPipePane(); err != nil {
		return err
	}

	if err := x.wait(); err != nil {
		return err
	}

	output, err := x.tmuxCapturePane()
	if err != nil {
		return err
	}
	c.Output = output

	if err := x.killSession(); err != nil && rerr == nil {
		rerr = err
	}
	sessionKilled = true

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

func (x *TmuxRunner) tmuxNewSession(args ...string) error {
	return x.tmux("new", append(
		[]string{"-s", x.session},
		args...,
	)...)
}

func (x *TmuxRunner) tmuxSendKeys(args ...string) error {
	return x.tmux("send-keys",
		"-t", x.session,
		strings.Join(args, " "),
		"C-m",
	)
}

func (x *TmuxRunner) startPipePane() error {
	// tmuxを直接呼び指すとうまく動かないので、シェル経由で呼び出す
	// cmd := exec.Command("tmux", "pipe-pane", "-t", "git-exec-session", "-o", "'cat >> /Users/akima/workspace/git-exec/trace.log'")
	cmd := exec.Command("/bin/zsh", "-c",
		fmt.Sprintf(
			"tmux pipe-pane -t git-exec-session -o 'cat >> %s'",
			x.pipeLogFile,
		),
	)

	slog.Debug("startPipePane", "args", cmd.Args)

	stderrBuf := new(strings.Builder)
	cmd.Stderr = stderrBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tmux pipe-pane: %w stderr: %s", err, stderrBuf.String())
	}
	return nil
}

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func (x *TmuxRunner) tmuxCapturePane() (string, error) {
	cmd := exec.Command("tmux", "capture-pane", "-t", x.session, "-pS", "-", "-e")
	slog.Debug("tmux capture-pane", "args", cmd.Args)
	b, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("tmux capture-pane: %w", err)
	}
	cleanedOutput := ansiEscape.ReplaceAll(b, []byte(""))
	return string(cleanedOutput), nil
}

func (x *TmuxRunner) killSession() error {
	return x.tmux("kill-session", "-t", x.session)
}

func (x *TmuxRunner) wait() error {
	for {
		time.Sleep(x.interval)

		found, err := x.findDoneString()
		if err != nil {
			return err
		}
		if found {
			break
		}
	}
	return nil
}

func (x *TmuxRunner) findDoneString() (bool, error) {
	logger := slog.Default().With("file", x.pipeLogFile)
	logger.Debug("findDoneString 0")
	tmpFile, err := os.Open(x.pipeLogFile)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Debug("findDoneString", "err", "not found")
			return false, nil
		}
		logger.Error("findDoneString", "err", err)
		return false, err
	}
	defer tmpFile.Close()

	logger.Debug("findDoneString 1")

	b, err := io.ReadAll(tmpFile)
	if err != nil {
		logger.Error("findDoneString", "read err", err)
		return false, err
	}
	logger.Debug("findDoneString 2", "length", len(b))
	if bytes.Contains(b, []byte(x.errorString)) {
		return false, fmt.Errorf("error occurred")
	}
	return bytes.Contains(b, []byte(x.completeString)), nil
}

func init() {
	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("failed to open debug.log, %w", err))
	}
	// defer f.Close()
	logger := slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
}
