package get_test

import (
	"testing"

	"github.com/chelnak/gh-changelog/internal/get"
	"github.com/stretchr/testify/assert"
)

var fileName string = "CHANGELOG.md"

func TestGetAll(t *testing.T) {
	cl, err := get.GetAll(fileName)

	// Should not error
	assert.Nil(t, err)

	// Should have at least 1 entry
	count := len(cl.GetEntries())
	assert.Greater(t, count, 0)
}

func TestGetLatest(t *testing.T) {
	cl, err := get.GetLatest(fileName)

	// Should not error
	assert.Nil(t, err)

	// Should have 1 entry
	count := len(cl.GetEntries())
	assert.Equal(t, 1, count)
}

func TestGetVersionWithAValidVersion(t *testing.T) {
	// Should not error when version is found
	cl, err := get.GetVersion(fileName, "v0.9.0")
	assert.Nil(t, err)

	// Should have 1 entry
	count := len(cl.GetEntries())
	assert.Equal(t, 1, count)
}

func TestGetVersionWithAnInvalidVersion(t *testing.T) {
	// Should error when version is not found
	_, err := get.GetVersion(fileName, "v0.0.0")
	assert.NotNil(t, err)
}
