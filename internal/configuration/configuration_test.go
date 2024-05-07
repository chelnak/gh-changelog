package configuration_test

import (
	"bytes"
	"testing"

	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/stretchr/testify/require"
)

func TestInitConfigSetsCorrectValues(t *testing.T) {
	err := configuration.InitConfig()
	require.NoError(t, err)

	config := configuration.Config

	require.Equal(t, "CHANGELOG.md", config.FileName)

	require.Equal(t, []string{"maintenance", "dependencies"}, config.ExcludedLabels)
	require.Equal(t, 2, len(config.ExcludedLabels))

	require.True(t, containsKey(config.Sections, "changed"))
	require.True(t, containsKey(config.Sections, "added"))
	require.True(t, containsKey(config.Sections, "fixed"))

	require.Equal(t, 3, len(config.Sections))
	require.Equal(t, false, config.SkipEntriesWithoutLabel)
	require.Equal(t, true, config.ShowUnreleased)
	require.Equal(t, true, config.CheckForUpdates)
	require.Equal(t, "spinner", config.Logger)
}

func TestPrintJSON(t *testing.T) {
	err := configuration.InitConfig()
	require.NoError(t, err)

	config := configuration.Config

	var buf bytes.Buffer
	err = config.PrintJSON(true, &buf)
	require.NoError(t, err)

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
  "checkForUpdates": true,
  "logger": "spinner"
}
`

	require.Equal(t, cfg, buf.String())
}

func TestPrintYAML(t *testing.T) {
	err := configuration.InitConfig()
	require.NoError(t, err)

	config := configuration.Config

	var buf bytes.Buffer
	err = config.PrintYAML(true, &buf)
	require.NoError(t, err)

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
logger: spinner
`
	require.Equal(t, cfg, buf.String())
}

func containsKey(m map[string][]string, key string) bool {
	_, ok := m[key]
	return ok
}
