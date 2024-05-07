package writer_test

import (
	"bytes"
	"regexp"
	"testing"
	"time"

	"github.com/chelnak/gh-changelog/internal/writer"
	"github.com/chelnak/gh-changelog/pkg/changelog"
	"github.com/chelnak/gh-changelog/pkg/entry"
	"github.com/stretchr/testify/require"
)

const (
	repoName  = "repo-name"
	repoOwner = "repo-owner"
)

func Test_ItWritesOutAChangelogInTheCorrectFormat(t *testing.T) {
	mockChangelog := changelog.NewChangelog(repoOwner, repoName)

	one := entry.Entry{
		Tag:        "v1.0.0",
		Date:       time.Now(),
		Added:      []string{"Added 1", "Added 2"},
		Changed:    []string{"Changed 1", "Changed 2"},
		Deprecated: []string{"Deprecated 1", "Deprecated 2"},
		Removed:    []string{"Removed 1", "Removed 2"},
		Fixed:      []string{"Fixed 1", "Fixed 2"},
		Security:   []string{"Security 1", "Security 2"},
		Other:      []string{"Other 1", "Other 2"},
	}

	two := one
	two.Tag = "v0.9.0"
	one.Previous = &two

	mockChangelog.Insert(one)
	mockChangelog.AddUnreleased([]string{"Unreleased 1", "Unreleased 2"})

	var buf bytes.Buffer
	err := writer.Write(&buf, mockChangelog)

	require.NoError(t, err)

	require.Regexp(t, "## Unreleased", buf.String())
	require.Regexp(t, "- Unreleased 1", buf.String())
	require.Regexp(t, "- Unreleased 2", buf.String())

	require.Regexp(t, regexp.MustCompile(`## \[v1.0.0]\(https://github.com/repo-owner/repo-name/tree/v1.0.0\)`), buf.String())
	require.Regexp(t, regexp.MustCompile(`\[Full Changelog]\(https://github.com/repo-owner/repo-name/compare/v0.9.0\.\.\.v1.0.0\)`), buf.String())

	require.Regexp(t, "### Added", buf.String())
	require.Regexp(t, "- Added 1", buf.String())
	require.Regexp(t, "- Added 2", buf.String())

	require.Regexp(t, "### Other", buf.String())
	require.Regexp(t, "- Other 1", buf.String())
	require.Regexp(t, "- Other 2", buf.String())
}
