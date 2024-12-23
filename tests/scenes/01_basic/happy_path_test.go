package basic

import (
	"testing"

	"github.com/akm/git-exec/core"
	"github.com/akm/git-exec/tests/testground"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHappyPath(t *testing.T) {
	defer testground.Setup(t)()

	// run(t, "make", "README.md")
	err := core.Run(core.DefaultOptions, []string{"make", "README.md"})
	require.NoError(t, err)

	commitMessage := stdout(t, "git", "log", "-1", "--pretty=%B")
	assert.Equal(t, `🤖 [01_basic] $ make README.md

Generating README.md

`, commitMessage)
}
