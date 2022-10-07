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
var fromVersion string
var fromLastVersion bool

// newCmd is the entry point for creating a new changelog
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new changelog from activity in the current repository",
	Long:  "Creates a new changelog from activity in the current repository.",
	RunE: func(command *cobra.Command, args []string) error {
		builder := changelog.NewChangelogBuilder()
		builder = builder.WithSpinner(true)
		builder = builder.WithNextVersion(nextVersion)
		builder = builder.WithFromVersion(fromVersion)
		builder = builder.WithFromLastVersion(fromLastVersion)

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

	newCmd.Flags().StringVar(
		&fromVersion,
		"from-version",
		"",
		"The version from which to start the changelog. If the value passed does not exist as a tag,\nthe changelog will be built from the first tag.",
	)

	newCmd.Flags().BoolVar(
		&fromLastVersion,
		"from-last-version",
		false,
		"Build the changelog starting from the last tag. Using this flag will result in a changelog with one entry.\nIt can be useful for generating a changelog to be used in release notes.",
	)

	newCmd.MarkFlagsMutuallyExclusive("from-version", "from-last-version")
	newCmd.Flags().SortFlags = false
}
