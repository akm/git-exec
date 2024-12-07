package main

import (
	"github.com/akm/git-exec/git"
	"github.com/akm/git-exec/opts"
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

const envKeyPrefix = "GIT_EXEC_"

func newOpt(shortName, longName string, hasValue bool, setFunc func(*Options, string)) *opts.Definition[Options] {
	return opts.NewDefinition(envKeyPrefix, shortName, longName, hasValue, setFunc)
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
	optDirectory = newOpt("-C", "--directory", true, func(o *Options, v string) { o.Directory = v }).WithoutEnv()
	optEmoji     = newOpt("-e", "--emoji", true, func(o *Options, v string) { o.Emoji = v })
	optPrompt    = newOpt("-p", "--prompt", true, func(o *Options, v string) { o.Prompt = v })
	optTemplate  = newOpt("-t", "--template", true, func(o *Options, v string) { o.Template = v })

	optSkipGuard                   = newOpt("", "--skip-guard", false, func(o *Options, _ string) { o.SkipGuard = true })
	optSkipGuardUncommittedChanges = newOpt("", "--skip-guard-uncommitted-changes", false, func(o *Options, _ string) { o.SkipGuardUncommittedChanges = true })
	optSkipGuardUntrackedFiles     = newOpt("", "--skip-guard-untracked-files", false, func(o *Options, _ string) { o.SkipGuardUntrackedFiles = true })

	optDebugLog    = newOpt("-D", "--debug-log", false, func(o *Options, _ string) { o.DebugLog = true })
	optInteractive = newOpt("-i", "--interactive", false, func(o *Options, _ string) { o.Interactive = true })

	optHelp    = newOpt("-h", "--help", false, func(o *Options, _ string) { o.Help = true }).WithoutEnv()
	optVersion = newOpt("-v", "--version", false, func(o *Options, _ string) { o.Version = true }).WithoutEnv()
)

var optionTypes = []*opts.Definition[Options]{
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

func newOptions() *Options {
	return opts.NewOptions(optionTypes, defaultOptions)
}

func parseOptions(args []string) (*Options, []string, error) {
	return opts.Parse(defaultOptions, optionTypes, args...)
}
