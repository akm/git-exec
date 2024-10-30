package main

import (
	"fmt"
	"strings"
)

type OptionType struct {
	ShortName string
	LongName  string
	HasValue  bool
}

func newOptionType(shortName, longName string, hasValue bool) *OptionType {
	return &OptionType{
		ShortName: shortName,
		LongName:  longName,
		HasValue:  hasValue,
	}
}

var (
	optHelp      = newOptionType("-h", "--help", false)
	optVersion   = newOptionType("-v", "--version", false)
	optDirectory = newOptionType("-C", "--directory", true)
)

var optionTypes = []*OptionType{
	optHelp,
	optVersion,
	optDirectory,
}

var optionKeyMap = func() map[string]*OptionType {
	m := map[string]*OptionType{}
	for _, opt := range optionTypes {
		m[opt.ShortName] = opt
		m[opt.LongName] = opt
	}
	return m
}()

type Option struct {
	Type  *OptionType
	Value string
}

type Options []*Option

func parseOptions(args []string) (Options, []string, error) {
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
