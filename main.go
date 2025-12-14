package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/akm/git-exec/core"
	"golang.org/x/exp/slog"
)

func main() {
	if len(os.Args) < 2 {
		core.Help()
		os.Exit(1)
	}

	logLevelMap := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	logLevel, ok := logLevelMap[strings.ToLower(os.Getenv("LOG_LEVEL"))]
	if !ok {
		logLevel = slog.LevelWarn
	}
	logHandler := slog.TextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(logHandler))

	options, commandArgs, err := core.ParseOptions(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse arguments: %s\n", err.Error())
	}
	if options.Help {
		core.Help()
		os.Exit(0)
	} else if options.Version {
		if len(commandArgs) == 0 {
			showVersion()
			os.Exit(0)
		} else {
			showVersionWithExecName(filepath.Base(os.Args[0]))
		}
	}

	if err := core.Run(options, commandArgs); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
