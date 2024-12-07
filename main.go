package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/akm/git-exec/git"
)

func main() {
	if len(os.Args) < 2 {
		help()
		os.Exit(1)
	}

	options, commandArgs, err := parseOptions(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse arguments: %s\n", err.Error())
	}
	if options.Help {
		help()
		os.Exit(0)
	} else if options.Version {
		if len(commandArgs) == 0 {
			showVersion()
			os.Exit(0)
		} else {
			showVersionWithExecName(filepath.Base(os.Args[0]))
		}
	}

	if err := process(options, commandArgs); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func process(options *Options, commandArgs []string) error {
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

	command := newCommand(commandArgs)
	var runner Runner
	if options.Interactive {
		runner = newTmuxRunner(options.DebugLog)
	} else {
		runner = newStandardRunner(options.DebugLog)
	}

	var commitMessage *commitMessage
	if err := changeDir((options.Directory), func() error {
		if err := runner.Run(command); err != nil {
			slog.Error("Command execution failed", "error", err)
			return fmt.Errorf("Command execution failed: %+v\n%s", err, command.Output)
		}
		commitMessage = newCommitMessage(command, options)
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

	if err := Commit(commitMessage); err != nil {
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
