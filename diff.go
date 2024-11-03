package main

import (
	"bytes"
	"os/exec"
)

func uncommittedChanges() (string, error) {
	output, err := exec.Command("git", "diff").CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(output)), nil
}

func hasUntrackedFiles() (bool, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(bytes.TrimSpace(output)) > 0, nil
}
