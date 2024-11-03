package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func commit(command *Command, options *Options) error {
	// 3. "git commit" を以下のオプションと標準力を指定して実行する。
	commitMessage, err := newCommitMessage(command, options).Build()
	if err != nil {
		return fmt.Errorf("Failed to build commit message: %+v", err)
	}

	// See https://tracpath.com/docs/git-commit/
	commitCmd := exec.Command("git", "commit", "--file", "-")
	commitCmd.Stdin = bytes.NewBufferString(commitMessage)

	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %+v", err)
	}

	return nil
}
