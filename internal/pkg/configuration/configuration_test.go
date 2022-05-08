package configuration_test

import (
	"testing"

	"github.com/chelnak/gh-changelog/internal/pkg/configuration"
	"github.com/stretchr/testify/assert"
)

func TestInitConfigSetsCorrectValues(t *testing.T) {
	err := configuration.InitConfig()
	assert.NoError(t, err)

	config := configuration.Config

	assert.Equal(t, "CHANGELOG.md", config.FileName)

	assert.Equal(t, []string{"maintenance"}, config.ExcludedLabels)
	assert.Equal(t, 1, len(config.ExcludedLabels))

	assert.True(t, containsKey(config.Sections, "Changed"))
	assert.True(t, containsKey(config.Sections, "Added"))
	assert.True(t, containsKey(config.Sections, "Fixed"))

	assert.Equal(t, 3, len(config.Sections))
	assert.Equal(t, false, config.SkipEntriesWithoutLabel)
	assert.Equal(t, true, config.ShowUnreleased)
	assert.Equal(t, true, config.CheckForUpdates)
}

func containsKey(m map[string][]string, key string) bool {
	_, ok := m[key]
	return ok
}
