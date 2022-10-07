// Package cmd holds all top-level cobra commands. Each file should contain
// only one command and that command should have only one purpose.
package cmd

import (
	"os"
	"path/filepath"

	"github.com/chelnak/gh-changelog/internal/changelog"
	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/writer"
	"github.com/spf13/cobra"
)

var nextVersion string

// newCmd is the entry point for creating a new changelog
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new changelog from activity in the current repository",
	Long:  "Creates a new changelog from activity in the current repository.",
	RunE: func(command *cobra.Command, args []string) error {
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
	},
}

func init() {
	newCmd.Flags().StringVar(&nextVersion, "next-version", "", "The next version to be released. The value passed does not have to be an existing tag.")
}
