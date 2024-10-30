package main

import (
	"fmt"
	"strings"
)

func splitStringsInto2(args []string, fn func(string) bool) ([]string, []string) {
	var a, b []string
	for _, arg := range args {
		if fn(arg) {
			a = append(a, arg)
		} else {
			b = append(b, arg)
		}
	}
	return a, b
}

func splitToOptionsAndCommandArgs(args []string) (Options, []string, error) {
	options := Options{}
	commandArgs := []string{}
	inOptions := true
	var waitingOption *Option
	for _, arg := range args {
		if waitingOption != nil {
			waitingOption.Value = arg
			options = append(options, waitingOption)
			waitingOption = nil
			continue
		}
		if inOptions && strings.HasPrefix(arg, "-") {
			optionType, ok := optionKeyMap[arg]
			if !ok {
				return nil, nil, fmt.Errorf("Unknown option: %s", arg)
			}
			if optionType.HasValue {
				waitingOption = &Option{Type: optionType}
			} else {
				options = append(options, &Option{Type: optionType})
			}
		} else {
			inOptions = false
			commandArgs = append(commandArgs, arg)
		}
	}
	if waitingOption != nil {
		return nil, nil, fmt.Errorf("no value given for option %s", waitingOption.Type.LongName)
	}
	return options, commandArgs, nil
}
