package changelog_test

import (
	"testing"
	"time"

	"github.com/chelnak/gh-changelog/pkg/changelog"
	"github.com/chelnak/gh-changelog/pkg/entry"
	"github.com/stretchr/testify/assert"
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

func TetstNewChangelog(t *testing.T) {
	var testChangelog = changelog.NewChangelog(repoName, repoOwner)
	assert.Equal(t, repoName, testChangelog.GetRepoName())
	assert.Equal(t, repoOwner, testChangelog.GetRepoOwner())
	assert.Equal(t, 0, len(testChangelog.GetEntries()))
	assert.Equal(t, 0, len(testChangelog.GetUnreleased()))
}

func TestInsert(t *testing.T) {
	var testChangelog = changelog.NewChangelog(repoName, repoOwner)
	for _, e := range entries {
		err := e.Append("added", "test")
		assert.Nil(t, err)

		testChangelog.Insert(e)
	}

	entries := testChangelog.GetEntries()
	assert.Equal(t, 2, len(entries))
	assert.Equal(t, 1, len(entries[0].Added))
	assert.Equal(t, "test", entries[0].Added[0])
}

func TestTail(t *testing.T) {
	var testChangelog = changelog.NewChangelog(repoName, repoOwner)

	for _, e := range entries {
		err := e.Append("added", "test")
		assert.Nil(t, err)

		testChangelog.Insert(e)
	}

	tail := testChangelog.Tail()
	assert.Equal(t, "v2.0.0", tail.Tag)
}

func TestHead(t *testing.T) {
	var testChangelog = changelog.NewChangelog(repoName, repoOwner)
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
		assert.Nil(t, err)

		testChangelog.Insert(e)
	}

	head := testChangelog.Head()
	assert.Equal(t, "v1.0.0", head.Tag)
}

func TestAddUnreleased(t *testing.T) {
	var testChangelog = changelog.NewChangelog(repoName, repoOwner)
	testChangelog.AddUnreleased([]string{"test"})

	unreleased := testChangelog.GetUnreleased()

	assert.Equal(t, 1, len(unreleased))
	assert.Equal(t, "test", unreleased[0])
}

func TestGetEntries(t *testing.T) {
	var testChangelog = changelog.NewChangelog(repoName, repoOwner)
	for _, e := range entries {
		err := e.Append("added", "test")
		assert.Nil(t, err)

		testChangelog.Insert(e)
	}

	assert.Equal(t, 2, len(testChangelog.GetEntries()))
	assert.Equal(t, "v2.0.0", testChangelog.GetEntries()[0].Tag)
}
