package basic

import (
	"testing"

	"github.com/akm/git-exec/core"
	"github.com/akm/git-exec/tests/testground"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectoryOption(t *testing.T) {
	defer testground.Setup(t)()

	defer testground.AssertStringNotChanged(t, testground.GitDiff)()

	opts := *core.DefaultOptions
	opts.Directory = "sub1"
	err := core.Run(&opts, []string{"make", "add-one", "parent-add-two"})
	require.NoError(t, err)

	commitMessage := stdout(t, "git", "log", "-1", "--pretty=%B")
	assert.Equal(t, `🤖 [02_subdir/sub1] $ make add-one parent-add-two

`, commitMessage)
}
