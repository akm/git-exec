package testground

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/akm/git-exec/testdir"
	"github.com/akm/git-exec/testexec"
)

func Setup(t *testing.T) func() {
	t.Helper()

	if os.Getenv("GITHUB_ACTIONS") == "true" {
		gitconfigPath := filepath.Join(os.Getenv("HOME"), ".gitconfig")
		gitConfigPath := filepath.Join(os.Getenv("HOME"), ".config", "git", "config")
		if !fileExists(gitconfigPath) && !fileExists(gitConfigPath) {
			testexec.Run(t, "git", "config", "--global", "user.email", "foo@example.com")
			testexec.Run(t, "git", "config", "--global", "user.name", "Foo Bar")
			testexec.Run(t, "git", "config", "--global", "init.defaultBranch", "main")
		}
	}

	// Suppress make's output
	os.Setenv("MAKEFLAGS", "--no-print-directory")

	r := testdir.Setup(t, ".", testdir.FromGoModRoot(t, "tests/grounds"))
	testexec.Run(t, "git", "init")
	testexec.Run(t, "git", "add", ".")
	testexec.Run(t, "git", "commit", "-m", "Initial commit")

	return r
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
