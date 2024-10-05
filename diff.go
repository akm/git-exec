package main

import (
	"os/exec"
)

func hasDiff() bool {
	if err := exec.Command("git", "diff", "--exit-code").Run(); err == nil {
		return false
	}
	return true
}
