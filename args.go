package main

import "strings"

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

func splitToOptionsAndCommandArgs(args []string) ([]string, []string) {
	inOptions := true
	return splitStringsInto2(args, func(arg string) bool {
		if inOptions && strings.HasPrefix(arg, "-") {
			return true
		} else {
			inOptions = false
			return false
		}
	})
}
