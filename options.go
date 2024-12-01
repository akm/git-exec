package main

import (
	"fmt"
	"os"
	"strings"
)

type Options struct {
	Help      bool
	Version   bool
	Directory string
	Emoji     string
	Prompt    string
	Template  string

	SkipGuard                   bool
	SkipGuardUncommittedChanges bool
	SkipGuardUntrackedFiles     bool

	DebugLog bool
	Tmux     bool
}

func newOptions() *Options {
	defaultOptionsCopy := *defaultOptions
	r := &defaultOptionsCopy
	for _, opt := range optionTypes {
		if opt.WithoutEnv {
			continue
		}
		if v := os.Getenv(opt.envKey()); v != "" {
			opt.setValue(r, v)
		}
	}
	return r
}

type OptionType struct {
	ShortName  string
	LongName   string
	HasValue   bool
	SetFunc    func(*Options, string)
	WithoutEnv bool
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

const envKeyPrefix = "GIT_EXEC_"

func (o *OptionType) envKey() string {
	return envKeyPrefix + strings.ToUpper(strings.ReplaceAll(strings.TrimLeft(o.LongName, "-"), "-", "_"))
}

func (o *OptionType) withoutEnv() *OptionType {
	o.WithoutEnv = true
	return o
}

var defaultOptions = &Options{
	Help:      false,
	Version:   false,
	Directory: "",
	Emoji:     "🤖",
	Prompt:    "$",
	Template:  `{{.Emoji}} [{{.Location}}] {{.Prompt}} {{.Command}}`,
	DebugLog:  false,
	Tmux:      false,
}

var (
	optDirectory = newOptionType("-C", "--directory", true, func(o *Options, v string) { o.Directory = v }).withoutEnv()
	optEmoji     = newOptionType("-e", "--emoji", true, func(o *Options, v string) { o.Emoji = v })
	optPrompt    = newOptionType("-p", "--prompt", true, func(o *Options, v string) { o.Prompt = v })
	optTemplate  = newOptionType("-t", "--template", true, func(o *Options, v string) { o.Template = v })

	optSkipGuard                   = newOptionType("", "--skip-guard", false, func(o *Options, _ string) { o.SkipGuard = true })
	optSkipGuardUncommittedChanges = newOptionType("", "--skip-guard-uncommitted-changes", false, func(o *Options, _ string) { o.SkipGuardUncommittedChanges = true })
	optSkipGuardUntrackedFiles     = newOptionType("", "--skip-guard-untracked-files", false, func(o *Options, _ string) { o.SkipGuardUntrackedFiles = true })

	optDebugLog = newOptionType("-D", "--debug-log", false, func(o *Options, _ string) { o.DebugLog = true })
	optTmux     = newOptionType("-T", "--tmux", false, func(o *Options, _ string) { o.Tmux = true })

	optHelp    = newOptionType("-h", "--help", false, func(o *Options, _ string) { o.Help = true }).withoutEnv()
	optVersion = newOptionType("-v", "--version", false, func(o *Options, _ string) { o.Version = true }).withoutEnv()
)

var optionTypes = []*OptionType{
	optDirectory,
	optEmoji,
	optPrompt,
	optTemplate,
	optSkipGuard,
	optSkipGuardUncommittedChanges,
	optSkipGuardUntrackedFiles,
	optDebugLog,
	optTmux,
	optHelp,
	optVersion,
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
	options := newOptions()
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
