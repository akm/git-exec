package basic

import (
	"testing"

	"github.com/akm/git-exec/core"
	"github.com/akm/git-exec/testexec"
	"github.com/akm/git-exec/tests/testground"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoChane(t *testing.T) {
	defer testground.Setup(t)()

	testexec.Run(t, "make", "README.md")
	testexec.Run(t, "git", "add", ".")
	testexec.Run(t, "git", "commit", "-m", "$ make README.md")

	lastCommitHash := testexec.Stdout(t, "git", "rev-parse", "HEAD")

	// testexec.Run(t, "make", "README.md")
	err := core.Run(core.DefaultOptions, []string{"make", "README.md"})
	require.Error(t, err)
	assert.Equal(t, "No changes to commit and No untracked files", err.Error())

	// No commit should be made
	currCommitHash := testexec.Stdout(t, "git", "rev-parse", "HEAD")
	assert.Equal(t, lastCommitHash, currCommitHash)
}