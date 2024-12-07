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

const envKeyPrefix = "GIT_EXEC_"

func boolOpt(shortName, longName string, setFunc func(*Options)) *opts.Definition[Options] {
	return opts.NewDefinition(envKeyPrefix, shortName, longName, false, func(o *Options, v string) { setFunc(o) })
}

func strOpt(shortName, longName string, setFunc func(*Options, string)) *opts.Definition[Options] {
	return opts.NewDefinition(envKeyPrefix, shortName, longName, true, setFunc)
}

var (
	optDirectory = strOpt("-C", "--directory", func(o *Options, v string) { o.Directory = v }).WithoutEnv()
	optEmoji     = strOpt("-e", "--emoji", func(o *Options, v string) { o.Emoji = v })
	optPrompt    = strOpt("-p", "--prompt", func(o *Options, v string) { o.Prompt = v })
	optTemplate  = strOpt("-t", "--template", func(o *Options, v string) { o.Template = v })

	optSkipGuard                   = boolOpt("", "--skip-guard", func(o *Options) { o.SkipGuard = true })
	optSkipGuardUncommittedChanges = boolOpt("", "--skip-guard-uncommitted-changes", func(o *Options) { o.SkipGuardUncommittedChanges = true })
	optSkipGuardUntrackedFiles     = boolOpt("", "--skip-guard-untracked-files", func(o *Options) { o.SkipGuardUntrackedFiles = true })

	optDebugLog    = boolOpt("-D", "--debug-log", func(o *Options) { o.DebugLog = true })
	optInteractive = boolOpt("-i", "--interactive", func(o *Options) { o.Interactive = true })

	optHelp    = boolOpt("-h", "--help", func(o *Options) { o.Help = true }).WithoutEnv()
	optVersion = boolOpt("-v", "--version", func(o *Options) { o.Version = true }).WithoutEnv()
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
