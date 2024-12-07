package main

import (
	"fmt"
	"strings"
)

func help() {
	firstLine := `Usage: git-exec [options ...] [key=value ...] <command> [args ...]`
	examples := `Examples:
* Specify environment variables.
	git exec FOO=fooooo make args1 args2

* Use shell to work with redirect operator.
	git exec /bin/bash -c 'echo "foo" >> README.md'

* Use interactive mode for command which requires input such as "npx sv create" for SvelteKit.
	git exec -i npx sv create my-app
`
	indent := "  "
	optionItems := make([]string, len(optionTypes))
	maxLongNameLength := 0
	for _, opt := range optionTypes {
		if maxLongNameLength < len(opt.LongName) {
			maxLongNameLength = len(opt.LongName)
		}
	}

	defaultOptions := newOptions()
	envVarItems := []string{}
	longNameFormat := "%-" + fmt.Sprintf("%ds", maxLongNameLength)
	for i, opt := range optionTypes {
		var item string
		if opt.ShortName == "" {
			item = fmt.Sprintf("%s    "+longNameFormat, indent, opt.LongName)
		} else {
			item = fmt.Sprintf("%s%s, "+longNameFormat, indent, opt.ShortName, opt.LongName)
		}
		item += " " + opt.GetHelp()
		if getter := opt.GetGetter(); getter != nil {
			item += fmt.Sprintf(" (default: %s)", getter(defaultOptions))
		}

		optionItems[i] = item
		if !opt.GetWithoutEnv() {
			envVarItems = append(envVarItems, fmt.Sprintf(longNameFormat+" %s", opt.LongName, opt.EnvKey()))
		}
	}
	options := "Options:\n" + strings.Join(optionItems, "\n")
	envVars := "Environment variable mapping:\n" + strings.Join(envVarItems, "\n")

	// git-exec は <command>よりも前に 複数のキーと値の組み合わせを指定可能で、
	// <command> 以後は 複数の引数を指定可能です。
	fmt.Println(firstLine + "\n\n" + options + "\n\n" + envVars + "\n\n" + examples)
}
