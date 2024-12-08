package basic

import (
	"testing"

	"github.com/akm/git-exec/core"
	"github.com/akm/git-exec/tests/testground"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommitMessage(t *testing.T) {
	defer testground.Setup(t)()

	opts := *core.DefaultOptions
	opts.Emoji = "🏭️"
	opts.Prompt = "%"
	opts.Template = `{{.Emoji}} {{.Prompt}} {{.Command}} [at {{.Location}}]`
	err := core.Run(&opts, []string{"make", "README.md"})
	require.NoError(t, err)

	commitMessage := stdout(t, "git", "log", "-1", "--pretty=%B")
	assert.Equal(t, `🏭️ % make README.md [at 01_basic]

Generating README.md

`, commitMessage)
}
