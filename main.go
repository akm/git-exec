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
	if options.Directory != "" {
		if err := os.Chdir(options.Directory); err != nil {
			return fmt.Errorf("Failed to change directory: %s", err.Error())
		}
	}

	if err := guard(options); err != nil {
		if isGuardError(err) {
			return err
		} else {
			return fmt.Errorf("Guard failed: %+v", err)
		}
	}

	command := newCommand(commandArgs)

	if err := command.Run(); err != nil {
		return fmt.Errorf("Command execution failed: %+v\n%s", err, command.Output)
	}

	if err := add(); err != nil {
		return err
	}

	if err := commit(command, options); err != nil {
		return err
	}

	return nil
}
