// Package cmd holds all top-level cobra commands. Each file should contain
// only one command and that command should have only one purpose.
package cmd

import (
	"fmt"

	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/pkg/parser"
	"github.com/spf13/cobra"
)

// parseCmd allows a user to parse a markdown changelog in to a struct.
// This is currently a work in progress and may go away in the future.
var parseCmd = &cobra.Command{
	Use:    "parse",
	Short:  "EXPERIMENTAL: Parse a changelog file in to a Changelog struct",
	Long:   "EXPERIMENTAL: Parse a changelog file in to a Changelog struct",
	Hidden: true,
	RunE: func(command *cobra.Command, args []string) error {
		changelog := configuration.Config.FileName

		parser := parser.NewParser(changelog, "", "")
		cl, err := parser.Parse()
		if err != nil {
			return err
		}

		// As an example, print out the changelog that was parsed.
		fmt.Printf("owner: %s\n", cl.GetRepoOwner())
		fmt.Printf("name: %s\n", cl.GetRepoName())
		fmt.Println("tags:")
		for _, e := range cl.GetEntries() {
			fmt.Printf(" %s\n", e.Tag)
		}

		return nil
	},
}
