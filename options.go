package main

import (
	"fmt"
	"strings"
)

type Options struct {
	Help      bool
	Version   bool
	Directory string
	Emoji     string
	Prompt    string
	Template  string
}

type OptionType struct {
	ShortName string
	LongName  string
	HasValue  bool
	SetFunc   func(*Options, string)
}

func newOptionType(shortName, longName string, hasValue bool, setFunc func(*Options, string)) *OptionType {
	return &OptionType{
		ShortName: shortName,
		LongName:  longName,
		HasValue:  hasValue,
		SetFunc:   setFunc,
	}
}

func (o *OptionType) setValue(opts *Options, value string) {
	o.SetFunc(opts, value)
}

var (
	optHelp      = newOptionType("-h", "--help", false, func(opts *Options, _ string) { opts.Help = true })
	optVersion   = newOptionType("-v", "--version", false, func(opts *Options, _ string) { opts.Version = true })
	optDirectory = newOptionType("-C", "--directory", true, func(opts *Options, value string) { opts.Directory = value })
	optEmoji     = newOptionType("-e", "--emoji", true, func(opts *Options, value string) { opts.Emoji = value })
	optPrompt    = newOptionType("-p", "--prompt", true, func(opts *Options, value string) { opts.Prompt = value })
	optTemplate  = newOptionType("-t", "--template", true, func(opts *Options, value string) { opts.Template = value })
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

func parseOptions(args []string) (*Options, []string, error) {
	options := &Options{}
	commandArgs := []string{}
	inOptions := true
	var pendingOptionType *OptionType
	for _, arg := range args {
		if pendingOptionType != nil {
			pendingOptionType.setValue(options, arg)
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
				optionType.setValue(options, "")
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
