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

var optionTypes = func() []*opts.Definition[Options] {
	envKeyPrefix := "GIT_EXEC_"
	boolOpt := func(shortName, longName string, setFunc func(*Options)) *opts.Definition[Options] {
		return opts.NewDefinition(envKeyPrefix, shortName, longName, false, func(o *Options, v string) { setFunc(o) })
	}
	strOpt := func(shortName, longName string, setFunc func(*Options, string)) *opts.Definition[Options] {
		return opts.NewDefinition(envKeyPrefix, shortName, longName, true, setFunc)
	}

	return []*opts.Definition[Options]{
		strOpt("-C", "--directory", func(o *Options, v string) { o.Directory = v }).WithoutEnv(),
		strOpt("-e", "--emoji", func(o *Options, v string) { o.Emoji = v }),
		strOpt("-p", "--prompt", func(o *Options, v string) { o.Prompt = v }),
		strOpt("-t", "--template", func(o *Options, v string) { o.Template = v }),

		boolOpt("", "--skip-guard", func(o *Options) { o.SkipGuard = true }),
		boolOpt("", "--skip-guard-uncommitted-changes", func(o *Options) { o.SkipGuardUncommittedChanges = true }),
		boolOpt("", "--skip-guard-untracked-files", func(o *Options) { o.SkipGuardUntrackedFiles = true }),

		boolOpt("-D", "--debug-log", func(o *Options) { o.DebugLog = true }),
		boolOpt("-i", "--interactive", func(o *Options) { o.Interactive = true }),

		boolOpt("-h", "--help", func(o *Options) { o.Help = true }).WithoutEnv(),
		boolOpt("-v", "--version", func(o *Options) { o.Version = true }).WithoutEnv(),
	}
}()

func newOptions() *Options {
	return opts.NewOptions(optionTypes, defaultOptions)
}

func parseOptions(args []string) (*Options, []string, error) {
	return opts.Parse(defaultOptions, optionTypes, args...)
}
