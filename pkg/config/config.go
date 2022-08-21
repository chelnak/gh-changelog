// Package config holds logic for working with application
// configuration.
package config

import (
	"errors"
	"io"

	"github.com/chelnak/gh-changelog/internal/configuration"
)

// Print prints the configuration to the terminal in either json or yaml format. Defaults to yaml.
func Print(output string, noColor bool, writer io.Writer) error {
	switch output {
	case "json":
		return configuration.Config.PrintJSON(noColor, writer)
	case "yaml":
		return configuration.Config.PrintYAML(noColor, writer)
	default:
		return errors.New("invalid output format. Valid values are 'json' and 'yaml'")
	}
}
