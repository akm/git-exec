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
	f := opts.NewFactory[Options]("GIT_EXEC_")

	return []*opts.Definition[Options]{
		f.String("-C", "--directory", "Specify the directory where the command is executed.",
			func(o *Options) string { return o.Directory },
			func(o *Options, v string) { o.Directory = v },
		).WithoutEnv(),
		f.String("-e", "--emoji", "Specify the emoji used in commit message.",
			func(o *Options) string { return o.Emoji },
			func(o *Options, v string) { o.Emoji = v },
		),
		f.String("-p", "--prompt", "Specify the prompt used in commit message.",
			func(o *Options) string { return o.Prompt },
			func(o *Options, v string) { o.Prompt = v },
		),
		f.String("-t", "--template", "Specify the template to build commit message.",
			func(o *Options) string { return o.Template },
			func(o *Options, v string) { o.Template = v },
		),

		f.Bool("", "--skip-guard", "Skip the guard check for uncommitted changes and untracked files before executing command.",
			func(o *Options) bool { return o.SkipGuard },
			func(o *Options) { o.SkipGuard = true },
		),
		f.Bool("", "--skip-guard-uncommitted-changes", "Skip the guard check for uncommitted changes before executing command.",
			func(o *Options) bool { return o.SkipGuardUncommittedChanges },
			func(o *Options) { o.SkipGuardUncommittedChanges = true },
		),
		f.Bool("", "--skip-guard-untracked-files", "Skip the guard check for untracked files before executing command.",
			func(o *Options) bool { return o.SkipGuardUntrackedFiles },
			func(o *Options) { o.SkipGuardUntrackedFiles = true },
		),

		f.Bool("-D", "--debug-log", "Output debug log.",
			func(o *Options) bool { return o.DebugLog },
			func(o *Options) { o.DebugLog = true },
		),
		f.Bool("-i", "--interactive", "Interactive mode for command which requires input. tmux is required to use.",
			func(o *Options) bool { return o.Interactive },
			func(o *Options) { o.Interactive = true },
		),

		f.Bool("-h", "--help", "Show this message.", nil, func(o *Options) { o.Help = true }).WithoutEnv(),
		f.Bool("-v", "--version", "Show version.", nil, func(o *Options) { o.Version = true }).WithoutEnv(),
	}
}()

func newOptions() *Options {
	return opts.NewOptions(optionTypes, defaultOptions)
}

func ParseOptions(args []string) (*Options, []string, error) {
	return opts.Parse(defaultOptions, optionTypes, args...)
}
