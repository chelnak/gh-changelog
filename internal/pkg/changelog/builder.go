package changelog

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chelnak/gh-changelog/internal/pkg/gitclient"
	"github.com/chelnak/gh-changelog/internal/pkg/githubclient"
	"github.com/chelnak/gh-changelog/internal/pkg/utils"
	"github.com/google/go-github/v43/github"
	"github.com/spf13/viper"
)

type Entry struct {
	Tag        string
	NextTag    string
	Date       time.Time
	Added      []string
	Changed    []string
	Deprecated []string
	Removed    []string
	Fixed      []string
	Security   []string
	Other      []string
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

func NewChangeLogBuilder(gitClient *gitclient.GitClient, githubClient *githubclient.GitHubClient, tags []*gitclient.Ref) *changeLogBuilder {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	_ = s.Color("green")
	s.FinalMSG = fmt.Sprintf("✅ Open %s or run 'gh changelog show' to view your changelog.\n", viper.GetString("file_name"))

	return &changeLogBuilder{
		spinner:      s,
		gitClient:    gitClient,
		githubClient: githubClient,
		tags:         tags,
	}
}

type changeLogBuilder struct {
	spinner      *spinner.Spinner
	gitClient    *gitclient.GitClient
	githubClient *githubclient.GitHubClient
	tags         []*gitclient.Ref
}

func (builder *changeLogBuilder) Build() (*ChangeLog, error) {
	changeLog := &ChangeLog{
		RepoName:  builder.githubClient.RepoContext.Name,
		RepoOwner: builder.githubClient.RepoContext.Owner,
		Entries:   []Entry{},
	}

	builder.spinner.Start()
	err := builder.buildChangeLog(changeLog)
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
		builder.spinner.Suffix = fmt.Sprintf(" Processing tags: 🏷️  %s", currentTag.Name)

		var nextTag *gitclient.Ref
		var err error
		if idx+1 == len(builder.tags) {
			nextTag, err = builder.gitClient.GetFirstCommit()
			if err != nil {
				return fmt.Errorf("could not get first commit: %v", err)
			}
		} else {
			nextTag = builder.tags[idx+1]
		}

		pullRequests, err := builder.githubClient.GetPullRequestsBetweenDates(nextTag.Date, currentTag.Date)
		if err != nil {
			return fmt.Errorf(
				"could not get pull requests for range '%s - %s': %v",
				nextTag.Date,
				currentTag.Date,
				err,
			)
		}

		entry, err := builder.populateEntry(
			currentTag.Name,
			nextTag.Name,
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

func (builder *changeLogBuilder) populateEntry(currentTag string, nextTag string, date time.Time, pullRequests []*github.Issue) (*Entry, error) {
	entry := &Entry{
		Tag:        currentTag,
		NextTag:    nextTag,
		Date:       date,
		Added:      []string{},
		Changed:    []string{},
		Deprecated: []string{},
		Removed:    []string{},
		Fixed:      []string{},
		Security:   []string{},
		Other:      []string{},
	}

	excludedLabels := viper.GetStringSlice("excluded_labels")
	for _, pr := range pullRequests {
		if !hasExcludedLabel(excludedLabels, pr) {
			line := fmt.Sprintf(
				"%s [#%d](https://github.com/%s/%s/pull/%d) ([%s](https://github.com/%s))\n",
				pr.GetTitle(),
				pr.GetNumber(),
				builder.githubClient.RepoContext.Owner,
				builder.githubClient.RepoContext.Name,
				pr.GetNumber(),
				pr.GetUser().GetLogin(),
				pr.GetUser().GetLogin(),
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

func hasExcludedLabel(excludedLabels []string, pr *github.Issue) bool {
	for _, label := range pr.Labels {
		if utils.Contains(excludedLabels, label.GetName()) {
			return true
		}
	}

	return false
}

func getSection(labels []*github.Label) string {
	sections := viper.GetStringMapStringSlice("sections")

	lookup := make(map[string]string)
	for k, v := range sections {
		for _, label := range v {
			lookup[label] = k
		}
	}

	section := ""
	skipUnlabelledEntries := viper.GetBool("skip_entries_without_label")
	for _, label := range labels {
		if _, ok := lookup[label.GetName()]; ok {
			section = lookup[label.GetName()]
		} else {
			if !skipUnlabelledEntries {
				section = "Other"
			}
		}
	}

	return section
}
