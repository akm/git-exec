package basic

import (
	"testing"

	"github.com/akm/git-exec/core"
	"github.com/akm/git-exec/tests/testground"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoChane(t *testing.T) {
	defer testground.Setup(t)()

	run(t, "make", "README.md")
	run(t, "git", "add", ".")
	run(t, "git", "commit", "-m", "$ make README.md")

	lastCommitHash := stdout(t, "git", "rev-parse", "HEAD")

	// run(t, "make", "README.md")
	err := core.Run(core.DefaultOptions, []string{"make", "README.md"})
	require.Error(t, err)
	assert.Equal(t, "No changes to commit and No untracked files", err.Error())

	// No commit should be made
	currCommitHash := stdout(t, "git", "rev-parse", "HEAD")
	assert.Equal(t, lastCommitHash, currCommitHash)
}
