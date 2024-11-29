package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	var origDir string
	if options.Directory != "" {
		{
			var err error
			origDir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("Failed to get current directory: %s", err.Error())
			}
		}
		if err := os.Chdir(options.Directory); err != nil {
			return fmt.Errorf("Failed to change directory: %s", err.Error())
		}
	}

	var guardMessage string
	if guardResult, err := guard(options); err != nil {
		return err
	} else if guardResult != nil {
		if guardResult.skipped {
			guardMessage = guardResult.Format()
			fmt.Fprintf(os.Stderr, "Guard skipped: %s\n", guardMessage)
		} else {
			return fmt.Errorf("Quit processing because %s", guardResult.Format())
		}
	}

	command := newCommand(commandArgs)

	if err := command.Run(); err != nil {
		return fmt.Errorf("Command execution failed: %+v\n%s", err, command.Output)
	}

	if origDir != "" {
		if err := os.Chdir(origDir); err != nil {
			return fmt.Errorf("Failed to change directory: %s", err.Error())
		}
	}

	if err := add(); err != nil {
		return err
	}

	commitMessage := newCommitMessage(command, options)

	if guardMessage != "" {
		commitMessage.Body = guardMessage + "\n\n" + commitMessage.Body
	}

	if err := commit(commitMessage); err != nil {
		return err
	}

	return nil
}
