package opts

import (
	"fmt"
	"strings"
)

type Definitions[T any] []*Definition[T]

func BuildKeyMap[T any](defs Definitions[T]) map[string]*Definition[T] {
	m := make(map[string]*Definition[T], len(defs))
	for _, def := range defs {
		m[def.LongName] = def
		m[def.ShortName] = def
	}
	return m
}

func Parse[T any](factory func() *T, defs Definitions[T], args ...string) (*T, []string, error) {
	options := factory()
	commandArgs := []string{}
	inOptions := true
	optionKeyMap := BuildKeyMap(defs)
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
