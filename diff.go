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

func hasUntrackedFiles() (bool, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(bytes.TrimSpace(output)) > 0, nil
}
