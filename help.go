package main

import "fmt"

func help() {
	// git-exec は <command>よりも前に 複数のキーと値の組み合わせを指定可能で、
	// <command> 以後は 複数の引数を指定可能です。
	fmt.Printf(`Usage: git-exec [key=value ...] <command> [args ...])

Examples:
* Specify environment variables.
	git-exec FOO=fooooo make args1 args2

* Use shell to work with redirect operator.
	git-exec /bin/bash -c 'echo "foo" >> README.md'
`)
}
