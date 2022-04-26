package changelog

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chelnak/gh-changelog/internal/pkg/gitclient"
	"github.com/chelnak/gh-changelog/internal/pkg/githubv4client"
	"github.com/chelnak/gh-changelog/internal/pkg/utils"
	"github.com/spf13/viper"
)

type Entry struct {
	CurrentTag  string
	PreviousTag string
	Date        time.Time
	Added       []string
	Changed     []string
	Deprecated  []string
	Removed     []string
	Fixed       []string
	Security    []string
	Other       []string
}

func (e *Entry) Append(section string, entry string) error {
	switch strings.ToLower(section) {
	case "added":
		e.Added = append(e.Added, entry)
	case "changed":
		e.Changed = append(e.Changed, entry)
	case "deprecated":
		e.Deprecated = append(e.Deprecated, entry)
	case "removed":
		e.Removed = append(e.Removed, entry)
	case "fixed":
		e.Fixed = append(e.Fixed, entry)
	case "security":
		e.Security = append(e.Security, entry)
	case "other":
		e.Other = append(e.Other, entry)
	default:
		return fmt.Errorf("unknown entry type '%s'", section)
	}

	return nil
}

type ChangeLog struct {
	RepoName  string
	RepoOwner string
	Entries   []Entry
}

type changeLogBuilder struct {
	spinner        *spinner.Spinner
	githubV4Client *githubv4client.GitHubGraphClient
	git            *gitclient.Git
	tags           []githubv4client.Tag
}

func (builder *changeLogBuilder) Build() (*ChangeLog, error) {
	builder.spinner.Start()

	builder.spinner.Suffix = " Fetching tags..."

	tags, err := builder.githubV4Client.GetTags()
	if err != nil {
		return nil, err
	}

	builder.tags = tags

	changeLog := &ChangeLog{
		RepoName:  builder.githubV4Client.RepoContext.Name,
		RepoOwner: builder.githubV4Client.RepoContext.Owner,
		Entries:   []Entry{},
	}

	err = builder.buildChangeLog(changeLog)
	if err != nil {
		builder.spinner.FinalMSG = ""
		builder.spinner.Stop()
		return nil, err
	}

	builder.spinner.Stop()
	return changeLog, nil
}

func (builder *changeLogBuilder) buildChangeLog(changeLog *ChangeLog) error {
	for idx, currentTag := range builder.tags {
		builder.spinner.Suffix = fmt.Sprintf(" Building changelog: 🏷️  %s", currentTag.Name)
		var previousTag githubv4client.Tag

		if idx+1 == len(builder.tags) {
			firstCommitSha, err := builder.git.GetFirstCommit()
			if err != nil {
				return err
			}

			date, err := builder.git.GetDateOfHash(firstCommitSha)
			if err != nil {
				return err
			}

			previousTag = githubv4client.Tag{
				Name: firstCommitSha,
				Sha:  firstCommitSha,
				Date: date,
			}
		} else {
			previousTag = builder.tags[idx+1]
		}

		pullRequests, err := builder.githubV4Client.GetPullRequestsBetweenDates(previousTag.Date, currentTag.Date)
		if err != nil {
			return err
		}

		entry, err := builder.populateEntry(
			currentTag.Name,
			previousTag.Name,
			currentTag.Date,
			pullRequests,
		)
		if err != nil {
			return fmt.Errorf("could not process pull requests: %v", err)
		}

		changeLog.Entries = append(changeLog.Entries, *entry)
	}

	return nil
}

func (builder *changeLogBuilder) populateEntry(currentTag string, previousTag string, date time.Time, pullRequests []githubv4client.PullRequest) (*Entry, error) {
	entry := &Entry{
		CurrentTag:  currentTag,
		PreviousTag: previousTag,
		Date:        date,
		Added:       []string{},
		Changed:     []string{},
		Deprecated:  []string{},
		Removed:     []string{},
		Fixed:       []string{},
		Security:    []string{},
		Other:       []string{},
	}

	excludedLabels := viper.GetStringSlice("excluded_labels")
	for _, pr := range pullRequests {
		if !hasExcludedLabel(excludedLabels, pr) {
			line := fmt.Sprintf(
				"%s [#%d](https://github.com/%s/%s/pull/%d) ([%s](https://github.com/%s))\n",
				pr.Title,
				pr.Number,
				builder.githubV4Client.RepoContext.Owner,
				builder.githubV4Client.RepoContext.Name,
				pr.Number,
				pr.User,
				pr.User,
			)

			section := getSection(pr.Labels)
			if section != "" {
				err := entry.Append(section, line)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return entry, nil
}

func NewChangeLogBuilder() (*changeLogBuilder, error) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	_ = s.Color("green")
	s.FinalMSG = fmt.Sprintf("✅ Open %s or run 'gh changelog show' to view your changelog.\n", viper.GetString("file_name"))

	v4, err := githubv4client.NewGitHubGraphClient()
	if err != nil {
		return nil, err
	}

	git := gitclient.NewGitHandler()

	builder := &changeLogBuilder{
		spinner:        s,
		githubV4Client: v4,
		git:            git,
		tags:           []githubv4client.Tag{},
	}

	return builder, err
}

func hasExcludedLabel(excludedLabels []string, pr githubv4client.PullRequest) bool {
	for _, label := range pr.Labels {
		if utils.Contains(excludedLabels, label.Name) {
			return true
		}
	}

	return false
}

func getSection(labels []githubv4client.Label) string {
	sections := viper.GetStringMapStringSlice("sections")

	lookup := make(map[string]string)
	for k, v := range sections {
		for _, label := range v {
			lookup[label] = k
		}
	}

	var section string
	skipUnlabelledEntries := viper.GetBool("skip_entries_without_label")

	if !skipUnlabelledEntries {
		section = "Other"
	}

	for _, label := range labels {
		if _, ok := lookup[label.Name]; ok {
			section = lookup[label.Name]
		}
	}

	return section
}
