package main

import (
	"bytes"
	"os/exec"
)

func hasUncommittedChanges() (bool, error) {
	output, err := exec.Command("git", "diff").CombinedOutput()
	if err != nil {
		return false, err
	}
	return len(bytes.TrimSpace(output)) > 0, nil
}
