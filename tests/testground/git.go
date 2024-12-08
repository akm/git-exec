package testground

import (
	"testing"

	"github.com/akm/git-exec/testexec"
)

func GitLastCommitHash(t *testing.T) string {
	t.Helper()
	return testexec.Stdout(t, "git", "rev-parse", "HEAD")
}

func GitDiff(t *testing.T) string {
	t.Helper()
	return testexec.Stdout(t, "git", "diff")
}
