package changelog_test

import (
	"testing"
	"time"

	"github.com/chelnak/gh-changelog/pkg/changelog"
	"github.com/stretchr/testify/assert"
)

const (
	repoName  = "repo-name"
	repoOwner = "repo-owner"
)

var testChangelog = changelog.NewChangelog(repoName, repoOwner)

func TetstNewChangelog(t *testing.T) {
	assert.Equal(t, repoName, testChangelog.GetRepoName())
	assert.Equal(t, repoOwner, testChangelog.GetRepoOwner())
	assert.Equal(t, 0, len(testChangelog.GetEntries()))
	assert.Equal(t, 0, len(testChangelog.GetUnreleased()))
}

func TestAddEntry(t *testing.T) {
	entry := changelog.Entry{
		CurrentTag:  "v2.0.0",
		PreviousTag: "v1.0.0",
		Date:        time.Time{},
	}

	err := entry.Append("added", "test")
	assert.Nil(t, err)

	testChangelog.AddEntry(entry)

	entries := testChangelog.GetEntries()

	assert.Equal(t, 1, len(entries))
	assert.Equal(t, 1, len(entries[0].Added))
	assert.Equal(t, "test", entries[0].Added[0])
}

func TestAddingAndRetrievingUnreleasedReturnsTheCorrectResponse(t *testing.T) {
	testChangelog.AddUnreleased([]string{"test"})

	unreleased := testChangelog.GetUnreleased()

	assert.Equal(t, 1, len(unreleased))
	assert.Equal(t, "test", unreleased[0])
}
