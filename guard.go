package main

import (
	"bytes"
	"os/exec"
)

type guardError struct {
	message string
}

func (e *guardError) Error() string {
	return e.message
}

func isGuardError(err error) bool {
	_, ok := err.(*guardError)
	return ok
}

func guard() error {
	// 環境変数 GIT_EXEC_SKIP_GUARD, GIT_EXEC_SKIP_GUARD_UNCOMMITED_CHANGES, GIT_EXEC_SKIP_GUARD_UNTRACKED_FILES には
	// 以下の値で真偽値を表す文字列を想定する
	// true: "true", "1", "yes", "on"
	// false: "false", "0", "no", "off"
	// 空文字列 あるいは それ以外の文字列は false として扱う
	//
	// GIT_EXEC_SKIP_GUARD あるいは GIT_EXEC_SKIP_GUARD_UNCOMMITED_CHANGES のいずれかが true でなければ、コミットされていない変更があればエラーを返す
	// GIT_EXEC_SKIP_GUARD あるいは GIT_EXEC_SKIP_GUARD_UNTRACKED_FILES のいずれかが true でなければ、追跡されていないファイルがあればエラーを返す

	if err := guardUncommitedFiles(); err != nil {
		return err
	}
	if err := guardUntrackedFiles(); err != nil {
		return err
	}

	return nil
}

func guardUncommitedFiles() error {
	if getEnvBool("GIT_EXEC_SKIP_GUARD") || getEnvBool("GIT_EXEC_SKIP_GUARD_UNCOMMITED_CHANGES") {
		return nil
	}
	if hasDiff(true) {
		return &guardError{"There are uncommitted changes"}
	}
	return nil
}

func guardUntrackedFiles() error {
	if getEnvBool("GIT_EXEC_SKIP_GUARD") || getEnvBool("GIT_EXEC_SKIP_GUARD_UNTRACKED_FILES") {
		return nil
	}
	r, err := hasUntrackedFiles()
	if err != nil {
		return err
	}
	if r {
		return &guardError{"There are untracked files"}
	}
	return nil
}

func hasUntrackedFiles() (bool, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(bytes.TrimSpace(output)) > 0, nil
}
