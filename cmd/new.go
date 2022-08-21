// Package cmd holds all top-level cobra commands. Each file should contain
// only one command and that command should have only one purpose.
package cmd

import (
	"github.com/chelnak/gh-changelog/pkg/new"
	"github.com/spf13/cobra"
)

var nextVersion string

// newCmd is the entry point for creating a new changelog
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new changelog from activity in the current repository",
	Long:  "Creates a new changelog from activity in the current repository.",
	RunE: func(command *cobra.Command, args []string) error {
		return new.CreateChangeLog(nextVersion)
	},
}

func init() {
	newCmd.Flags().StringVar(&nextVersion, "next-version", "", "The next version to be released. The value passed does not have to be an existing tag.")
}
