package core

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/akm/git-exec/command"
	"github.com/akm/git-exec/git"
)

func Run(options *Options, commandArgs []string) error {
	slog.Debug("Run started", "options", options, "commandArgs", commandArgs)

	var guardMessage string
	if guardResult, err := git.Guard(&options.GuardOptions); err != nil {
		return err
	} else if guardResult != nil {
		if guardResult.Skipped {
			guardMessage = guardResult.Format()
			fmt.Fprintf(os.Stderr, "Guard skipped: %s\n", guardMessage)
		} else {
			return fmt.Errorf("Quit processing because %s", guardResult.Format())
		}
	}

	cmd := command.NewCommand(commandArgs)
	var runner command.Runner
	if options.Interactive {
		runner = command.NewTmuxRunner(options.DebugLog)
	} else {
		runner = command.NewStandardRunner(options.DebugLog)
	}

	var commitMessage *commitMessage
	if err := changeDir((options.Directory), func() error {
		if err := runner.Run(cmd); err != nil {
			slog.Error("Command execution failed", "error", err)
			return fmt.Errorf("Command execution failed: %+v\n%s", err, cmd.Output)
		}
		commitMessage = newCommitMessage(cmd, options)

		location, err := getLocation()
		if err != nil {
			return err
		}
		commitMessage.Location = location

		return nil
	}); err != nil {
		return err
	}

	if err := git.Add(); err != nil {
		return err
	}

	if guardMessage != "" {
		commitMessage.Body = guardMessage + "\n\n" + commitMessage.Body
	}

	// 3. "git commit" を以下のオプションと標準力を指定して実行する。
	msg, err := commitMessage.Build()
	if err != nil {
		return fmt.Errorf("Failed to build commit message: %+v", err)
	}

	if err := git.Commit(msg); err != nil {
		return err
	}

	return nil
}

func changeDir(dir string, cb func() error) (rerr error) {
	if dir == "" {
		return cb()
	}
	var origDir string
	if dir != "" {
		{
			var err error
			origDir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("Failed to get current directory: %s", err.Error())
			}
		}
		if err := os.Chdir(dir); err != nil {
			return fmt.Errorf("Failed to change directory: %s", err.Error())
		}
	}
	if origDir != "" {
		defer func() {
			if err := os.Chdir(origDir); err != nil {
				rerr = fmt.Errorf("Failed to change directory: %s", err.Error())
			}
		}()
	}
	return cb()
}
