package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/akm/git-exec/git"
)

type Options struct {
	Help      bool
	Version   bool
	Directory string
	Emoji     string
	Prompt    string
	Template  string

	git.GuardOptions

	DebugLog    bool
	Interactive bool
}

func newOptions() *Options {
	defaultOptionsCopy := *defaultOptions
	r := &defaultOptionsCopy
	for _, opt := range optionTypes {
		if opt.WithoutEnv {
			continue
		}
		if v := os.Getenv(opt.envKey()); v != "" {
			opt.SetFunc(r, v)
		}
	}
	return r
}

const envKeyPrefix = "GIT_EXEC_"

func newOpt(shortName, longName string, hasValue bool, setFunc func(*Options, string)) *OptionType {
	return newOptionType(envKeyPrefix, shortName, longName, hasValue, setFunc)
}

var defaultOptions = &Options{
	Help:        false,
	Version:     false,
	Directory:   "",
	Emoji:       "🤖",
	Prompt:      "$",
	Template:    `{{.Emoji}} [{{.Location}}] {{.Prompt}} {{.Command}}`,
	DebugLog:    false,
	Interactive: false,
}

var (
	optDirectory = newOpt("-C", "--directory", true, func(o *Options, v string) { o.Directory = v }).withoutEnv()
	optEmoji     = newOpt("-e", "--emoji", true, func(o *Options, v string) { o.Emoji = v })
	optPrompt    = newOpt("-p", "--prompt", true, func(o *Options, v string) { o.Prompt = v })
	optTemplate  = newOpt("-t", "--template", true, func(o *Options, v string) { o.Template = v })

	optSkipGuard                   = newOpt("", "--skip-guard", false, func(o *Options, _ string) { o.SkipGuard = true })
	optSkipGuardUncommittedChanges = newOpt("", "--skip-guard-uncommitted-changes", false, func(o *Options, _ string) { o.SkipGuardUncommittedChanges = true })
	optSkipGuardUntrackedFiles     = newOpt("", "--skip-guard-untracked-files", false, func(o *Options, _ string) { o.SkipGuardUntrackedFiles = true })

	optDebugLog    = newOpt("-D", "--debug-log", false, func(o *Options, _ string) { o.DebugLog = true })
	optInteractive = newOpt("-i", "--interactive", false, func(o *Options, _ string) { o.Interactive = true })

	optHelp    = newOpt("-h", "--help", false, func(o *Options, _ string) { o.Help = true }).withoutEnv()
	optVersion = newOpt("-v", "--version", false, func(o *Options, _ string) { o.Version = true }).withoutEnv()
)

var optionTypes = OptionTypes{
	optDirectory,
	optEmoji,
	optPrompt,
	optTemplate,
	optSkipGuard,
	optSkipGuardUncommittedChanges,
	optSkipGuardUntrackedFiles,
	optDebugLog,
	optInteractive,
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
