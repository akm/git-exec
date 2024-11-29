package main

import "fmt"

const Version = "0.0.14"

func showVersion() {
	fmt.Println(Version)
}

func showVersionWithExecName(execName string) {
	fmt.Printf("%s@v%s\n", execName, Version)
}
