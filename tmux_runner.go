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
	session     string
	doneString  string
	pipeLogFile string
	interval    time.Duration
	debugLog    bool
}

var _ Runner = (*TmuxRunner)(nil)

func newTmuxRunner(debugLog bool) *TmuxRunner {
	return &TmuxRunner{
		session:     "git-exec-session",
		doneString:  "git-exec-done",
		pipeLogFile: "/Users/akima/workspace/git-exec/trace.log",
		debugLog:    debugLog,
		interval:    1_000 * time.Millisecond,
	}
}

func (x *TmuxRunner) Run(c *Command) (rerr error) {
	if _, err := exec.LookPath("tmux"); err != nil {
		return fmt.Errorf("tmux is not installed. Please install tmux. See https://github.com/tmux/tmux/wiki/Installing")
	}

	tracingFile := "/Users/akima/workspace/git-exec/trace.log"

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

	inputs := append(append(c.Envs, c.Args...), "; echo "+x.doneString)
	if err := x.tmuxSendKeys(inputs...); err != nil {
		return err
	}

	// tmpFile, err := os.CreateTemp("", "git-exec-pipe-pane")
	// if err != nil {
	// 	return err
	// }
	// if err := tmpFile.Close(); err != nil {
	// 	return err
	// }
	// defer os.Remove(tmpFile.Name())

	if err := x.startPipePane(); err != nil {
		return err
	}

	for {
		time.Sleep(x.interval)

		found, err := x.findDoneStringFromPipePane(tracingFile)
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
	arguments := []string{
		"-t",
		x.session,
		strings.Join(args, " "),
		"C-m",
	}
	return x.tmux("send-keys", arguments...)
}

func (x *TmuxRunner) startPipePane() error {
	// cmd := exec.Command("tmux", append(
	// 	[]string{
	// 		"pipe-pane", "-t",
	// 		x.session,
	// 	},
	// 	args...,
	// )...)

	// cmd := exec.Command(os.Getenv("SHELL"), "-c",
	// 	singleQuote("tmux pipe-pane -t "+x.session+" "+strings.Join(args, " ")),
	// )

	// tmuxを直接呼び指すとうまく動かないので、シェル経由で呼び出す
	// cmd := exec.Command("tmux", "pipe-pane", "-t", "git-exec-session", "-o", "'cat >> /Users/akima/workspace/git-exec/trace.log'")
	cmd := exec.Command("/bin/zsh", "-c",
		fmt.Sprintf(
			"tmux pipe-pane -t git-exec-session -o 'cat >> %s'",
			x.pipeLogFile,
		),
	)
	// cmd.WaitDelay = 1 * time.Second

	slog.Debug("startPipePane", "args", cmd.Args)

	stderrBuf := new(strings.Builder)
	cmd.Stderr = stderrBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tmux pipe-pane: %w stderr: %s", err, stderrBuf.String())
	}
	return nil
}

func (x *TmuxRunner) stopPipePane() error {
	cmd := exec.Command("/bin/zsh", "-c", "tmux pipe-pane -t git-exec-session")
	// cmd.WaitDelay = 1 * time.Second

	slog.Debug("stopPipePane", "args", cmd.Args)

	stderrBuf := new(strings.Builder)
	cmd.Stderr = stderrBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tmux pipe-pane: %w stderr: %s", err, stderrBuf.String())
	}
	return nil
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

func doubleQuote(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}

func (x *TmuxRunner) findDoneStringFromPipePane(tmpFileName string) (bool, error) {
	slog.Debug("findDoneStringFromPipePane 0", "file", tmpFileName)
	tmpFile, err := os.Open(tmpFileName)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Debug("findDoneStringFromPipePane", "file", tmpFileName, "err", "not found")
			return false, nil
		}
		slog.Error("findDoneStringFromPipePane", "file", tmpFileName, "err", err)
		return false, err
	}
	defer tmpFile.Close()

	slog.Debug("findDoneStringFromPipePane 1", "file", tmpFileName)

	b, err := io.ReadAll(tmpFile)
	if err != nil {
		slog.Error("findDoneStringFromPipePane", "file", tmpFileName, "read err", err)
		return false, err
	}
	slog.Debug("findDoneStringFromPipePane 2", "file", tmpFileName, "length", len(b))
	return bytes.Contains(b, []byte(x.doneString)), nil
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
