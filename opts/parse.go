package opts

import (
	"fmt"
	"strings"
)

func Parse[T any](factory func() *T, defs []*Definition[T], args ...string) (*T, []string, error) {
	options := factory()
	commandArgs := []string{}
	inOptions := true
	optionKeyMap := buildKeyMap(defs)
	var pendingOptionType *Definition[T]
	for _, arg := range args {
		if pendingOptionType != nil {
			pendingOptionType.SetFunc(options, arg)
			pendingOptionType = nil
			continue
		}
		if inOptions && strings.HasPrefix(arg, "-") {
			optionType, ok := optionKeyMap[arg]
			if !ok {
				return nil, nil, fmt.Errorf("Unknown option: %s", arg)
			}
			if optionType.HasValue {
				pendingOptionType = optionType
			} else {
				optionType.SetFunc(options, "")
			}
		} else {
			inOptions = false
			commandArgs = append(commandArgs, arg)
		}
	}
	if pendingOptionType != nil {
		return nil, nil, fmt.Errorf("no value given for option %s", pendingOptionType.LongName)
	}
	return options, commandArgs, nil
}

func buildKeyMap[T any](defs []*Definition[T]) map[string]*Definition[T] {
	m := make(map[string]*Definition[T], len(defs))
	for _, def := range defs {
		m[def.LongName] = def
		m[def.ShortName] = def
	}
	return m
}
