package main

import "strings"

func splitToOptionsAndCommandArgs(args []string) ([]string, []string) {
	var options []string
	var commandArgs []string
	inOptions := true
	for _, arg := range args {
		if inOptions && strings.HasPrefix(arg, "-") {
			options = append(options, arg)
		} else {
			inOptions = false
			commandArgs = append(commandArgs, arg)
		}
	}
	return options, commandArgs
}
