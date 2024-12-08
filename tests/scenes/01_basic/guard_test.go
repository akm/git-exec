package basic

import (
	"strings"
	"testing"

	"github.com/akm/git-exec/core"
	"github.com/akm/git-exec/testdir"
	"github.com/akm/git-exec/testexec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuardUntrackedFiles(t *testing.T) {
	defer testdir.Setup(t, ".", testdir.FromGitRoot(t, "tests/grounds"))()
	testexec.Run(t, "git", "init")
	testexec.Run(t, "git", "add", ".")
	testexec.Run(t, "git", "commit", "-m", "Initial commit")
	lastCommitHash := testexec.Stdout(t, "git", "rev-parse", "HEAD")

	testexec.Run(t, "make", "add-one") // Let it not be committed

	err := core.Run(core.DefaultOptions, []string{"make", "README.md"})
	require.Error(t, err)
	assert.Equal(t, `Quit processing because There are untracked files

Untracked files:
work.txt`, err.Error())

	// No commit should be made
	currCommitHash := testexec.Stdout(t, "git", "rev-parse", "HEAD")
	assert.Equal(t, lastCommitHash, currCommitHash)
}

func TestGuardUncommittedChanes(t *testing.T) {
	defer testdir.Setup(t, ".", testdir.FromGitRoot(t, "tests/grounds"))()
	testexec.Run(t, "git", "init")
	testexec.Run(t, "git", "add", ".")
	testexec.Run(t, "git", "commit", "-m", "Initial commit")

	// commit add-one
	testexec.Run(t, "make", "add-one")
	testexec.Run(t, "git", "add", ".")
	testexec.Run(t, "git", "commit", "-m", "add one")

	lastCommitHash := testexec.Stdout(t, "git", "rev-parse", "HEAD")

	testexec.Run(t, "make", "add-two") // Let it not be committed

	err := core.Run(core.DefaultOptions, []string{"make", "README.md"})
	require.Error(t, err)
	assert.Equal(t, `Quit processing because There are uncommitted changes

Uncommitted changes:
`+strings.TrimSpace(testexec.Stdout(t, "git", "diff")),
		err.Error())

	// No commit should be made
	currCommitHash := testexec.Stdout(t, "git", "rev-parse", "HEAD")
	assert.Equal(t, lastCommitHash, currCommitHash)
}
