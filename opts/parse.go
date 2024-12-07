package opts

import (
	"fmt"
	"os"
	"strings"
)

func Parse[T any](defualtOptions *T, defs []*Definition[T], args ...string) (*T, []string, error) {
	options := NewOptions(defs, defualtOptions)
	commandArgs := []string{}
	inOptions := true
	optionKeyMap := buildKeyMap(defs)
	var pendingOptionType *Definition[T]
	for _, arg := range args {
		if pendingOptionType != nil {
			pendingOptionType.Setter(options, arg)
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
				optionType.Setter(options, "")
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

func NewOptions[T any](defs []*Definition[T], defaultOptions *T) *T {
	copy := *defaultOptions
	r := &copy
	for _, opt := range defs {
		if opt.GetWithoutEnv() {
			continue
		}
		if v := os.Getenv(opt.EnvKey()); v != "" {
			opt.Setter(r, v)
		}
	}
	return r
}

func buildKeyMap[T any](defs []*Definition[T]) map[string]*Definition[T] {
	m := make(map[string]*Definition[T], len(defs))
	for _, def := range defs {
		m[def.LongName] = def
		m[def.ShortName] = def
	}
	return m
}
