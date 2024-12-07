package main

import (
	"fmt"
	"strings"

	"github.com/akm/git-exec/opts"
)

func Help() {
	firstLine := `Usage: git-exec [options ...] [key=value ...] <command> [args ...]`
	examples := `Examples:
* Specify environment variables.
	git exec FOO=fooooo make args1 args2

* Use shell to work with redirect operator.
	git exec /bin/bash -c 'echo "foo" >> README.md'

* Use interactive mode for command which requires input such as "npx sv create" for SvelteKit.
	git exec -i npx sv create my-app
`
	optionItems, envVarItems := opts.HelpItemsAndEnvVarMappings[Options](defaultOptions, optionTypes)

	options := "Options:\n" + strings.Join(optionItems, "\n")
	envVars := "Environment variable mapping:\n" + strings.Join(envVarItems, "\n")

	// git-exec は <command>よりも前に 複数のキーと値の組み合わせを指定可能で、
	// <command> 以後は 複数の引数を指定可能です。
	fmt.Println(firstLine + "\n\n" + options + "\n\n" + envVars + "\n\n" + examples)
}
