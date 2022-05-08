package cmd

import (
	"github.com/chelnak/gh-changelog/internal/pkg/configuration"
	"github.com/chelnak/gh-changelog/internal/pkg/markdown"
	"github.com/spf13/cobra"
)

// showCmd is the entry point for rendering a changelog in the terminal
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Renders the current changelog in the terminal",
	Long:  "Renders the current changelog in the terminal",
	RunE: func(command *cobra.Command, args []string) error {
		changelog := configuration.Config.FileName
		return markdown.Render(changelog)
	},
}
