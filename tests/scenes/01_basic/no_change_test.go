package basic

import (
	"testing"

	"github.com/akm/git-exec/core"
	"github.com/akm/git-exec/tests/testground"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoChange(t *testing.T) {
	defer testground.Setup(t)()

	run(t, "make", "README.md")
	run(t, "git", "add", ".")
	run(t, "git", "commit", "-m", "$ make README.md")

	defer testground.AssertStringNotChanged(t, testground.GitLastCommitHash)()

	err := core.Run(core.DefaultOptions, []string{"make", "README.md"})
	require.Error(t, err)
	assert.Equal(t, "No changes to commit and No untracked files", err.Error())
}
