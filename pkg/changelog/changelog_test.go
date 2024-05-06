package changelog_test

import (
	"testing"
	"time"

	"github.com/chelnak/gh-changelog/pkg/changelog"
	"github.com/chelnak/gh-changelog/pkg/entry"
	"github.com/stretchr/testify/require"
)

const (
	repoName  = "repo-name"
	repoOwner = "repo-owner"
)

var entries = []entry.Entry{
	{
		Tag:  "v2.0.0",
		Date: time.Time{},
	},
	{
		Tag:  "v1.0.0",
		Date: time.Time{},
	},
}

func TestNewChangelog(t *testing.T) {
	var testChangelog = changelog.NewChangelog(repoOwner, repoName)
	require.Equal(t, repoName, testChangelog.GetRepoName())
	require.Equal(t, repoOwner, testChangelog.GetRepoOwner())
	require.Equal(t, 0, len(testChangelog.GetEntries()))
	require.Equal(t, 0, len(testChangelog.GetUnreleased()))
}

func TestInsert(t *testing.T) {
	var testChangelog = changelog.NewChangelog(repoOwner, repoName)
	for _, e := range entries {
		err := e.Append("added", "test")
		require.Nil(t, err)

		testChangelog.Insert(e)
	}

	entries := testChangelog.GetEntries()
	require.Equal(t, 2, len(entries))
	require.Equal(t, 1, len(entries[0].Added))
	require.Equal(t, "test", entries[0].Added[0])
}

func TestTail(t *testing.T) {
	var testChangelog = changelog.NewChangelog(repoOwner, repoName)

	for _, e := range entries {
		err := e.Append("added", "test")
		require.Nil(t, err)

		testChangelog.Insert(e)
	}

	tail := testChangelog.Tail()
	require.Equal(t, "v2.0.0", tail.Tag)
}

func TestHead(t *testing.T) {
	var testChangelog = changelog.NewChangelog(repoOwner, repoName)
	entries := []entry.Entry{
		{
			Tag:  "v2.0.0",
			Date: time.Time{},
		},
		{
			Tag:  "v1.0.0",
			Date: time.Time{},
		},
	}

	for _, e := range entries {
		err := e.Append("added", "test")
		require.Nil(t, err)

		testChangelog.Insert(e)
	}

	head := testChangelog.Head()
	require.Equal(t, "v1.0.0", head.Tag)
}

func TestAddUnreleased(t *testing.T) {
	var testChangelog = changelog.NewChangelog(repoOwner, repoName)
	testChangelog.AddUnreleased([]string{"test"})

	unreleased := testChangelog.GetUnreleased()

	require.Equal(t, 1, len(unreleased))
	require.Equal(t, "test", unreleased[0])
}

func TestGetEntries(t *testing.T) {
	var testChangelog = changelog.NewChangelog(repoOwner, repoName)
	for _, e := range entries {
		err := e.Append("added", "test")
		require.Nil(t, err)

		testChangelog.Insert(e)
	}

	require.Equal(t, 2, len(testChangelog.GetEntries()))
	require.Equal(t, "v2.0.0", testChangelog.GetEntries()[0].Tag)
}
