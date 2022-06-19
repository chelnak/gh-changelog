package configuration_test

import (
	"bytes"
	"testing"

	"github.com/chelnak/gh-changelog/internal/pkg/configuration"
	"github.com/stretchr/testify/assert"
)

func TestInitConfigSetsCorrectValues(t *testing.T) {
	err := configuration.InitConfig()
	assert.NoError(t, err)

	config := configuration.Config

	assert.Equal(t, "CHANGELOG.md", config.FileName)

	assert.Equal(t, []string{"maintenance", "dependencies"}, config.ExcludedLabels)
	assert.Equal(t, 2, len(config.ExcludedLabels))

	assert.True(t, containsKey(config.Sections, "changed"))
	assert.True(t, containsKey(config.Sections, "added"))
	assert.True(t, containsKey(config.Sections, "fixed"))

	assert.Equal(t, 3, len(config.Sections))
	assert.Equal(t, false, config.SkipEntriesWithoutLabel)
	assert.Equal(t, true, config.ShowUnreleased)
	assert.Equal(t, true, config.CheckForUpdates)
}

func TestPrintJSON(t *testing.T) {
	err := configuration.InitConfig()
	assert.NoError(t, err)

	config := configuration.Config

	var buf bytes.Buffer
	err = config.PrintJSON(true, &buf)
	assert.NoError(t, err)

	cfg := `{
  "fileName": "CHANGELOG.md",
  "excludedLabels": [
    "maintenance",
    "dependencies"
  ],
  "sections": {
    "added": [
      "feature",
      "enhancement"
    ],
    "changed": [
      "backwards-incompatible"
    ],
    "fixed": [
      "bug",
      "bugfix",
      "documentation"
    ]
  },
  "skipEntriesWithoutLabel": false,
  "showUnreleased": true,
  "checkForUpdates": true
}
`
	assert.Equal(t, cfg, buf.String())
}

func TestPrintYAML(t *testing.T) {
	err := configuration.InitConfig()
	assert.NoError(t, err)

	config := configuration.Config

	var buf bytes.Buffer
	err = config.PrintYAML(true, &buf)
	assert.NoError(t, err)

	cfg := `---
file_name: CHANGELOG.md
excluded_labels:
- maintenance
- dependencies
sections:
  added:
  - feature
  - enhancement
  changed:
  - backwards-incompatible
  fixed:
  - bug
  - bugfix
  - documentation
skip_entries_without_label: false
show_unreleased: true
check_for_updates: true
`
	assert.Equal(t, cfg, buf.String())
}

func containsKey(m map[string][]string, key string) bool {
	_, ok := m[key]
	return ok
}
