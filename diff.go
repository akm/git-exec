package main

import (
	"fmt"
	"os/exec"
)

func hasDiff(showOutput bool) bool {
	output, err := exec.Command("git", "diff", "--exit-code").CombinedOutput()
	if err != nil {
		if showOutput {
			fmt.Println(string(output))
		}
		return true
	}
	return false
}
