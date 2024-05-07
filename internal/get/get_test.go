package get_test

import (
	"testing"

	"github.com/chelnak/gh-changelog/internal/get"
	"github.com/stretchr/testify/require"
)

var fileName = "CHANGELOG.md"
var singleEntryFileName = "single_CHANGELOG.md"

func TestGetAll(t *testing.T) {
	cl, err := get.GetAll(fileName)

	// Should not error
	require.Nil(t, err)

	// Should have at least 1 entry
	count := len(cl.GetEntries())
	require.Greater(t, count, 0)
}

func TestGetLatest(t *testing.T) {
	cl, err := get.GetLatest(fileName)

	// Should not error
	require.Nil(t, err)

	// Should have 1 entry
	count := len(cl.GetEntries())
	require.Equal(t, 1, count)
	require.Equal(t, "v0.13.0", cl.GetEntries()[0].PrevTag)
}

func TestGetLatestWithNoPrevious(t *testing.T) {
	cl, err := get.GetLatest(singleEntryFileName)

	// Should not error
	require.Nil(t, err)

	// Should have 1 entry
	count := len(cl.GetEntries())
	require.Equal(t, 1, count)
	require.Equal(t, "", cl.GetEntries()[0].PrevTag)
}

func TestGetVersionWithAValidVersion(t *testing.T) {
	// Should not error when version is found
	cl, err := get.GetVersion(fileName, "v0.9.0")
	require.Nil(t, err)

	// Should have 1 entry
	count := len(cl.GetEntries())
	require.Equal(t, 1, count)
	require.Equal(t, "v0.8.1", cl.GetEntries()[0].PrevTag)
}

func TestGetVersionWithAnInvalidVersion(t *testing.T) {
	// Should error when version is not found
	_, err := get.GetVersion(fileName, "v0.0.0")
	require.NotNil(t, err)
}
