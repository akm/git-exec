package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Commit(msg string) error {
	// See https://tracpath.com/docs/git-commit/
	commitCmd := exec.Command("git", "commit", "--file", "-")
	commitCmd.Stdin = bytes.NewBufferString(msg)

	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %+v", err)
	}

	return nil
}
