package entry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/chelnak/gh-changelog/pkg/changelog"
	"github.com/chelnak/gh-changelog/pkg/entry"
	"github.com/stretchr/testify/require"
)

const (
	repoName  = "gh-changelog"
	repoOwner = "chelnak"
)

func TestNewEntry(t *testing.T) {
	e := entry.NewEntry("v2.0.0", time.Time{})

	require.Equal(t, "v2.0.0", e.Tag)
	require.Equal(t, time.Time{}, e.Date)
}

func TestPrevious(t *testing.T) {
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

	tail := testChangelog.Tail()
	previous := tail.Previous
	require.Equal(t, "v1.0.0", previous.Tag)
}

func TestNext(t *testing.T) {
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
	next := head.Next
	require.Equal(t, "v2.0.0", next.Tag)
}

func TestAppend(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "added",
		},
		{
			name: "changed",
		},
		{
			name: "deprecated",
		},
		{
			name: "removed",
		},
		{
			name: "fixed",
		},
		{
			name: "security",
		},
		{
			name: "other",
		},
	}

	e := entry.NewEntry("v2.0.0", time.Time{})
	for _, test := range tests {
		t.Run(fmt.Sprintf("Appends a line to section: %s", test.name), func(t *testing.T) {
			err := e.Append(test.name, fmt.Sprintf("test %s", test.name))
			require.Nil(t, err)

			section := e.GetSection(test.name)
			require.Equal(t, 1, len(section))
			require.Regexp(t, fmt.Sprintf("test %s", test.name), section[0])
		})
	}
}

func TestReturnsAnErrorWhenAppendingToAnInvalidSection(t *testing.T) {
	e := entry.NewEntry("v2.0.0", time.Time{})
	err := e.Append("invalid", "test")
	require.NotNil(t, err)
}

func TestGetSection(t *testing.T) {
	e := entry.NewEntry("v2.0.0", time.Time{})
	err := e.Append("added", "test")
	require.Nil(t, err)

	section := e.GetSection("added")
	require.Equal(t, 1, len(section))
	require.Equal(t, "test", section[0])

	section = e.GetSection("invalid")
	require.Equal(t, 0, len(section))
}
