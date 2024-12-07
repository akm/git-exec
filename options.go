package main

import (
	"fmt"

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
	boolToString := func(b bool) string {
		if b {
			return "true"
		} else {
			return "false"
		}
	}
	quote := func(s string) string {
		return fmt.Sprintf("%q", s)
	}

	envKeyPrefix := "GIT_EXEC_"
	boolOpt := func(shortName, longName string, help string, getter func(*Options) bool, setter func(*Options)) *opts.Definition[Options] {
		var actualGetter func(o *Options) string
		if getter != nil {
			actualGetter = func(o *Options) string { return boolToString(getter(o)) }
		}
		return opts.NewDefinition(envKeyPrefix, shortName, longName, false,
			func(o *Options, v string) { setter(o) }).
			Getter(actualGetter).Help(help)
	}
	strOpt := func(shortName, longName string, help string, getter func(*Options) string, setter func(*Options, string)) *opts.Definition[Options] {
		var actualGetter func(o *Options) string
		if getter != nil {
			actualGetter = func(o *Options) string { return quote(getter(o)) }
		}
		return opts.NewDefinition(envKeyPrefix, shortName, longName, true, setter).
			Getter(actualGetter).Help(help)
	}

	return []*opts.Definition[Options]{
		strOpt("-C", "--directory", "Specify the directory where the command is executed.",
			func(o *Options) string { return o.Directory },
			func(o *Options, v string) { o.Directory = v },
		).WithoutEnv(),
		strOpt("-e", "--emoji", "Specify the emoji used in commit message.",
			func(o *Options) string { return o.Emoji },
			func(o *Options, v string) { o.Emoji = v },
		),
		strOpt("-p", "--prompt", "Specify the prompt used in commit message.",
			func(o *Options) string { return o.Prompt },
			func(o *Options, v string) { o.Prompt = v },
		),
		strOpt("-t", "--template", "Specify the template to build commit message.",
			func(o *Options) string { return o.Template },
			func(o *Options, v string) { o.Template = v },
		),

		boolOpt("", "--skip-guard", "Skip the guard check for uncommitted changes and untracked files before executing command.",
			func(o *Options) bool { return o.SkipGuard },
			func(o *Options) { o.SkipGuard = true },
		),
		boolOpt("", "--skip-guard-uncommitted-changes", "Skip the guard check for uncommitted changes before executing command.",
			func(o *Options) bool { return o.SkipGuardUncommittedChanges },
			func(o *Options) { o.SkipGuardUncommittedChanges = true },
		),
		boolOpt("", "--skip-guard-untracked-files", "Skip the guard check for untracked files before executing command.",
			func(o *Options) bool { return o.SkipGuardUntrackedFiles },
			func(o *Options) { o.SkipGuardUntrackedFiles = true },
		),

		boolOpt("-D", "--debug-log", "Output debug log.",
			func(o *Options) bool { return o.DebugLog },
			func(o *Options) { o.DebugLog = true },
		),
		boolOpt("-i", "--interactive", "Interactive mode for command which requires input. tmux is required to use.",
			func(o *Options) bool { return o.Interactive },
			func(o *Options) { o.Interactive = true },
		),

		boolOpt("-h", "--help", "Show this message.", nil, func(o *Options) { o.Help = true }).WithoutEnv(),
		boolOpt("-v", "--version", "Show version.", nil, func(o *Options) { o.Version = true }).WithoutEnv(),
	}
}()

func newOptions() *Options {
	return opts.NewOptions(optionTypes, defaultOptions)
}

func parseOptions(args []string) (*Options, []string, error) {
	return opts.Parse(defaultOptions, optionTypes, args...)
}
