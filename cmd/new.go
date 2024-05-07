// Package cmd holds all top-level cobra commands. Each file should contain
// only one command and that command should have only one purpose.
package cmd

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/log"
	"github.com/chelnak/gh-changelog/internal/writer"
	"github.com/chelnak/gh-changelog/pkg/builder"
	"github.com/spf13/cobra"
)

var (
	nextVersion   string
	fromVersion   string
	latestVersion bool
	filter        string
	ancestorsOnly bool
	logger        string
)

// newCmd is the entry point for creating a new changelog
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new changelog from activity in the current repository",
	Long:  "Creates a new changelog from activity in the current repository.",
	RunE: func(command *cobra.Command, args []string) error {
		log.SetupLogging(log.GetLoggerType(logger))
		builder, err := builder.NewBuilder()
		if err != nil {
			return err
		}

		if nextVersion != "" {
			builder.NextVersion(nextVersion)
		}

		if fromVersion != "" {
			builder.FromVersion(fromVersion)
		}

		if latestVersion {
			builder.LatestVersion()
		}

		if filter != "" {
			fil, err := regexp.Compile(filter)
			if err != nil {
				return err
			}
			builder.Filter(fil)
		}

		if ancestorsOnly {
			builder.AncestorsOnly()
		}

		changelog, err := builder.Build()
		if err != nil {
			return err
		}

		f, err := os.Create(filepath.Clean(configuration.Config.FileName))
		if err != nil {
			return err
		}

		if err := writer.Write(f, changelog); err != nil {
			return err
		}

		return nil
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
		&latestVersion,
		"latest",
		false,
		"Build the changelog starting from the latest tag. Using this flag will result in a changelog with one entry.\nIt can be useful for generating a changelog to be used in release notes.",
	)

	newCmd.Flags().StringVar(&filter, "filter", "", "Filter the results by tag name. This flag supports regular expressions.")
	newCmd.Flags().BoolVar(&ancestorsOnly, "ancestors-only", false, "Builds the changelog with tags that are ancestor of the current branch.")
	newCmd.Flags().StringVar(&logger, "logger", "spinner", "The type of logger to use. Valid values are 'spinner' and 'console'. The default is 'spinner'.")

	newCmd.MarkFlagsMutuallyExclusive("from-version", "latest")
	newCmd.Flags().SortFlags = false
}
