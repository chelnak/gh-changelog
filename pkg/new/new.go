// Package new holds logic for the creation of the projects changelog.
package new

import (
	"os"
	"path/filepath"

	"github.com/chelnak/gh-changelog/internal/changelog"
	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/writer"
)

// CreateChangeLog is the entry point for creating a changelog in the termina	.
func CreateChangeLog(nextVersion string) error {
	builder := changelog.NewChangelogBuilder()
	builder = builder.WithSpinner(true)
	builder = builder.WithNextVersion(nextVersion)

	changelog, err := builder.Build()
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Clean(configuration.Config.FileName))
	if err != nil {
		return err
	}

	return writer.Write(f, changelog)
}
