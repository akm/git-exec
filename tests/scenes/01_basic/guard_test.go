package basic

import (
	"strings"
	"testing"

	"github.com/akm/git-exec/core"
	"github.com/akm/git-exec/tests/testground"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuardUntrackedFiles(t *testing.T) {
	defer testground.Setup(t)()
	lastCommitHash := stdout(t, "git", "rev-parse", "HEAD")

	run(t, "make", "add-one") // Let it not be committed

	err := core.Run(core.DefaultOptions, []string{"make", "README.md"})
	require.Error(t, err)
	assert.Equal(t, `Quit processing because There are untracked files

Untracked files:
work.txt`, err.Error())

	// No commit should be made
	currCommitHash := stdout(t, "git", "rev-parse", "HEAD")
	assert.Equal(t, lastCommitHash, currCommitHash)
}

func TestGuardUncommittedChanes(t *testing.T) {
	defer testground.Setup(t)()

	// commit add-one
	run(t, "make", "add-one")
	run(t, "git", "add", ".")
	run(t, "git", "commit", "-m", "add one")

	lastCommitHash := stdout(t, "git", "rev-parse", "HEAD")

	run(t, "make", "add-two") // Let it not be committed

	err := core.Run(core.DefaultOptions, []string{"make", "README.md"})
	require.Error(t, err)
	assert.Equal(t, `Quit processing because There are uncommitted changes

Uncommitted changes:
`+strings.TrimSpace(stdout(t, "git", "diff")),
		err.Error())

	// No commit should be made
	currCommitHash := stdout(t, "git", "rev-parse", "HEAD")
	assert.Equal(t, lastCommitHash, currCommitHash)
}
