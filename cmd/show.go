package cmd

import (
	"github.com/chelnak/gh-changelog/internal/pkg/markdown"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// showCmd is the entry point for rendering a changelog in the terminal
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Renders the current changelog in the terminal",
	Long:  "Renders the current changelog in the terminal",
	RunE: func(command *cobra.Command, args []string) error {
		changelog := viper.GetString("file_name")
		return markdown.Render(changelog)
	},
}
