package cmd

import (
	"github.com/chelnak/gh-changelog/internal/pkg/changelog"
	"github.com/chelnak/gh-changelog/internal/pkg/writer"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new changelog from activity in the current repository",
	Long:  "Creates a new changelog from activity the current repository.",
	RunE: func(command *cobra.Command, args []string) error {
		changeLog, err := changelog.NewChangelog()
		if err != nil {
			return err
		}
		return writer.Write(changeLog)
	},
}
