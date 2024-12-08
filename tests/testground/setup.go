package testground

import (
	"os"
	"testing"

	"github.com/akm/git-exec/testdir"
	"github.com/akm/git-exec/testexec"
)

func Setup(t *testing.T) func() {
	t.Helper()

	if os.Getenv("GITHUB_ACTIONS") == "true" {
		// git config --global user.email "foo@example.com"
		// git config --global user.name "Foo Bar"
		// git config --global init.defaultBranch main
		testexec.Run(t, "git", "config", "--global", "user.email", "foo@example.com")
		testexec.Run(t, "git", "config", "--global", "user.name", "Foo Bar")
		testexec.Run(t, "git", "config", "--global", "init.defaultBranch", "main")
	}

	r := testdir.Setup(t, ".", testdir.FromGoModRoot(t, "tests/grounds"))
	testexec.Run(t, "git", "init")
	testexec.Run(t, "git", "add", ".")
	testexec.Run(t, "git", "commit", "-m", "Initial commit")

	return r
}
