package main

import (
	"fmt"
	"os/exec"
)

func add() error {
	uncommittedChanges, err := hasUncommittedChanges()
	if err != nil {
		return fmt.Errorf("git diff failed: %+v", err)
	}
	untrackedFiles, err := hasUntrackedFiles()
	if err != nil {
		return fmt.Errorf("git ls-files failed: %+v", err)
	}

	if !uncommittedChanges && !untrackedFiles {
		return fmt.Errorf("No changes to commit and No untracked files")
	}

	if err := exec.Command("git", "add", ".").Run(); err != nil {
		return fmt.Errorf("git add failed: %+v", err)
	}

	return nil
}
