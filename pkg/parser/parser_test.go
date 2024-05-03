package parser_test

import (
	"testing"

	"github.com/chelnak/gh-changelog/pkg/parser"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Run("can parse a changelog with an unreleased section", func(t *testing.T) {
		p := parser.NewParser("./testdata/unreleased.md", "chelnak", "gh-changelog")
		c, err := p.Parse()
		require.NoError(t, err)
		require.Len(t, c.GetUnreleased(), 2)
		require.Len(t, c.GetEntries(), 3)
	})

	t.Run("can parse a changelog without an unreleased section", func(t *testing.T) {
		p := parser.NewParser("./testdata/no_unreleased.md", "chelnak", "gh-changelog")
		c, err := p.Parse()
		require.NoError(t, err)
		require.Len(t, c.GetUnreleased(), 0)
		require.Len(t, c.GetEntries(), 3)
	})
}
