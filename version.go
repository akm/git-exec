package main

import "fmt"

const Version = "0.1.1"

func showVersion() {
	fmt.Println(Version)
}

func showVersionWithExecName(execName string) {
	fmt.Printf("%s@v%s\n", execName, Version)
}
