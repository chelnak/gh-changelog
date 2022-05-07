package changelog

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chelnak/gh-changelog/internal/pkg/gitclient"
	"github.com/chelnak/gh-changelog/internal/pkg/githubclient"
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

func (e *Entry) append(section string, entry string) error {
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

type Changelog interface {
	GetRepoName() string
	GetRepoOwner() string
	GetEntries() []Entry
}

type changelog struct {
	repoName  string
	repoOwner string
	entries   []Entry
}

type ChangelogBuilder interface {
	WithSpinner(enabled bool) ChangelogBuilder
	WithGitClient(client gitclient.GitClient) ChangelogBuilder
	WithGithubClient(client githubclient.GitHubClient) ChangelogBuilder
	Build() (Changelog, error)
}

type changelogBuilder struct {
	spinner *spinner.Spinner
	github  githubclient.GitHubClient
	git     gitclient.GitClient
	tags    []githubclient.Tag
}

func (builder *changelogBuilder) WithSpinner(enabled bool) ChangelogBuilder {
	if enabled {
		builder.spinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		_ = builder.spinner.Color("green")
		builder.spinner.FinalMSG = fmt.Sprintf("✅ Open %s or run 'gh changelog show' to view your changelog.\n", viper.GetString("file_name"))
	}
	return builder
}

func (builder *changelogBuilder) WithGitClient(git gitclient.GitClient) ChangelogBuilder {
	builder.git = git
	return builder
}

func (builder *changelogBuilder) WithGithubClient(client githubclient.GitHubClient) ChangelogBuilder {
	builder.github = client
	return builder
}

func (builder *changelogBuilder) Build() (Changelog, error) {
	if builder.spinner != nil {
		builder.spinner.Start()
	}

	var err error
	if builder.github == nil {
		builder.github, err = githubclient.NewGitHubClient()
		if err != nil {
			return nil, err
		}
	}

	if builder.git == nil {
		builder.git = gitclient.NewGitClient()
	}

	builder.spinner.Suffix = " Fetching tags..."
	tags, err := builder.github.GetTags()
	if err != nil {
		return nil, err
	}

	builder.tags = tags

	c := &changelog{
		repoName:  builder.github.GetRepoName(),
		repoOwner: builder.github.GetRepoOwner(),
		entries:   []Entry{},
	}

	err = builder.buildChangeLog(c)
	if err != nil {
		builder.spinner.FinalMSG = ""
		return nil, err
	}

	builder.spinner.Stop()

	return c, nil
}

func (builder *changelogBuilder) buildChangeLog(changelog *changelog) error {
	for idx, currentTag := range builder.tags {
		builder.spinner.Suffix = fmt.Sprintf(" Building changelog: 🏷️  %s", currentTag.Name)
		var previousTag githubclient.Tag

		if idx+1 == len(builder.tags) {
			firstCommitSha, err := builder.git.GetFirstCommit()
			if err != nil {
				return err
			}

			date, err := builder.git.GetDateOfHash(firstCommitSha)
			if err != nil {
				return err
			}

			previousTag = githubclient.Tag{
				Name: firstCommitSha,
				Sha:  firstCommitSha,
				Date: date,
			}
		} else {
			previousTag = builder.tags[idx+1]
		}

		pullRequests, err := builder.github.GetPullRequestsBetweenDates(previousTag.Date, currentTag.Date)
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

		changelog.entries = append(changelog.entries, *entry)
	}

	return nil
}

func (builder *changelogBuilder) populateEntry(currentTag string, previousTag string, date time.Time, pullRequests []githubclient.PullRequest) (*Entry, error) {
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
				builder.github.GetRepoOwner(),
				builder.github.GetRepoName(),
				pr.Number,
				pr.User,
				pr.User,
			)

			section := getSection(pr.Labels)
			if section != "" {
				err := entry.append(section, line)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return entry, nil
}

func (c *changelog) GetRepoName() string {
	return c.repoName
}

func (c *changelog) GetRepoOwner() string {
	return c.repoOwner
}

func (c *changelog) GetEntries() []Entry {
	return c.entries
}

func NewChangelogBuilder() ChangelogBuilder {
	return &changelogBuilder{}
}

func hasExcludedLabel(excludedLabels []string, pr githubclient.PullRequest) bool {
	for _, label := range pr.Labels {
		if utils.SliceContainsString(excludedLabels, label.Name) {
			return true
		}
	}

	return false
}

func getSection(labels []githubclient.Label) string {
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
