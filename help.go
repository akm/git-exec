package main

import (
	"fmt"
	"strings"
)

func help() {
	firstLine := `Usage: git-exec [options ...] [key=value ...] <command> [args ...]`
	examples := `Examples:
* Specify environment variables.
	git-exec FOO=fooooo make args1 args2

* Use shell to work with redirect operator.
	git-exec /bin/bash -c 'echo "foo" >> README.md'
`
	indent := "  "
	optionItems := make([]string, len(optionTypes))
	maxLongNameLength := 0
	for _, opt := range optionTypes {
		if maxLongNameLength < len(opt.LongName) {
			maxLongNameLength = len(opt.LongName)
		}
	}

	envVarItems := []string{}
	longNameFormat := "%-" + fmt.Sprintf("%ds", maxLongNameLength)
	for i, opt := range optionTypes {
		var item string
		if opt.ShortName == "" {
			item = fmt.Sprintf("%s    "+longNameFormat, indent, opt.LongName)
		} else {
			item = fmt.Sprintf("%s%s, "+longNameFormat, indent, opt.ShortName, opt.LongName)
		}
		item += " " + optionMessageMap[opt.LongName]
		optionItems[i] = item
		if !opt.WithoutEnv {
			envVarItems = append(envVarItems, fmt.Sprintf(longNameFormat+" %s", opt.LongName, opt.envKey()))
		}
	}
	options := "Options:\n" + strings.Join(optionItems, "\n")
	envVars := "Environment variable mapping:\n" + strings.Join(envVarItems, "\n")

	// git-exec は <command>よりも前に 複数のキーと値の組み合わせを指定可能で、
	// <command> 以後は 複数の引数を指定可能です。
	fmt.Println(firstLine + "\n\n" + options + "\n\n" + envVars + "\n\n" + examples)
}

var optionMessageMap = map[string]string{
	"--help":                           "Show this message.",
	"--version":                        "Show version.",
	"--directory":                      "Specify the directory where the command is executed.",
	"--emoji":                          "Specify the emoji used in commit message.",
	"--prompt":                         "Specify the prompt used in commit message.",
	"--template":                       "Specify the template to build commit message.",
	"--skip-guard":                     "Skip the guard check for uncommitted changes and untracked files before executing command.",
	"--skip-guard-uncommitted-changes": "Skip the guard check for uncommitted changes before executing command.",
	"--skip-guard-untracked-files":     "Skip the guard check for untracked files before executing command.",
}
