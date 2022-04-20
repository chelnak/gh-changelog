package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chelnak/gh-changelog/internal/pkg/changelog"
	"github.com/chelnak/gh-changelog/internal/pkg/writer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new changelog from activity in the current repository",
	Long:  "Creates a new changelog from activity the current repository.",
	RunE: func(command *cobra.Command, args []string) error {
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		_ = s.Color("green")
		s.FinalMSG = fmt.Sprintf("✅ Open %s or run 'gh changelog show' to view your changelog.\n", viper.GetString("file_name"))

		changeLog, err := changelog.MakeFullChangelog(s)
		if err != nil {
			return err
		}

		s.Stop()
		return writer.Write(changeLog)
	},
}
