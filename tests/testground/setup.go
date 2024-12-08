package testground

import (
	"testing"

	"github.com/akm/git-exec/testdir"
	"github.com/akm/git-exec/testexec"
)

func Setup(t *testing.T) func() {
	t.Helper()

	r := testdir.Setup(t, ".", testdir.FromGoModRoot(t, "tests/grounds"))
	testexec.Run(t, "git", "init")
	testexec.Run(t, "git", "add", ".")
	testexec.Run(t, "git", "commit", "-m", "Initial commit")

	return r
}
