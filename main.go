package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		Help()
		os.Exit(1)
	}

	options, commandArgs, err := ParseOptions(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse arguments: %s\n", err.Error())
	}
	if options.Help {
		Help()
		os.Exit(0)
	} else if options.Version {
		if len(commandArgs) == 0 {
			showVersion()
			os.Exit(0)
		} else {
			showVersionWithExecName(filepath.Base(os.Args[0]))
		}
	}

	if err := Run(options, commandArgs); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
