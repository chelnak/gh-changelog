// Package changelog package is responsible for collecting the data that
// will be used to generate the changelog.
package changelog

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/gitclient"
	"github.com/chelnak/gh-changelog/internal/githubclient"
	"github.com/chelnak/gh-changelog/internal/utils"
)

var Now = time.Now // must be a better way to stub this

// ChangelogBuilder represents
type ChangelogBuilder interface {
	WithSpinner(enabled bool) ChangelogBuilder
	WithGitClient(client gitclient.GitClient) ChangelogBuilder
	WithGitHubClient(client githubclient.GitHubClient) ChangelogBuilder
	WithNextVersion(nextVersion string) ChangelogBuilder
	Build() (Changelog, error)
}

type changelogBuilder struct {
	spinner     *spinner.Spinner
	github      githubclient.GitHubClient
	git         gitclient.GitClient
	tags        []githubclient.Tag
	nextVersion string
}

// NewChangelogBuilder creates a new returns a new instance of the changelog builder struct
func NewChangelogBuilder() ChangelogBuilder {
	return &changelogBuilder{}
}

// WithSpinner enables or disables the spinner for the changelog builder
func (builder *changelogBuilder) WithSpinner(enabled bool) ChangelogBuilder {
	if enabled {
		builder.spinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		_ = builder.spinner.Color("green")
		builder.spinner.FinalMSG = fmt.Sprintf("✅ Open %s or run 'gh changelog show' to view your changelog.\n", configuration.Config.FileName)
	}
	return builder
}

// WithGitClient allows the consumer to use a custom git client that implements the GitClient interface
func (builder *changelogBuilder) WithGitClient(git gitclient.GitClient) ChangelogBuilder {
	builder.git = git
	return builder
}

// WithGitHubClient allows the consumer to use a custom github client that implements the GitHubClient interface
func (builder *changelogBuilder) WithGitHubClient(client githubclient.GitHubClient) ChangelogBuilder {
	builder.github = client
	return builder
}

// WithNextVersion sets the next version to be used in the changelog. The value is either an empty string
// or a valid semantic version string passed from the consumer.
func (builder *changelogBuilder) WithNextVersion(nextVersion string) ChangelogBuilder {
	builder.nextVersion = nextVersion
	return builder
}

// Build builds the struct that is used to generate the changelog
func (builder *changelogBuilder) Build() (Changelog, error) {
	if builder.spinner != nil {
		builder.spinner.Start()
		defer builder.spinner.Stop()
	}

	var err error
	if builder.github == nil {
		builder.github, err = githubclient.NewGitHubClient()
		if err != nil {
			builder.spinner.FinalMSG = ""
			return nil, err
		}
	}

	if builder.git == nil {
		builder.git = gitclient.NewGitClient(exec.Command)
	}

	builder.spinner.Suffix = " Fetching tags..."
	tags, err := builder.github.GetTags()
	if err != nil {
		builder.spinner.FinalMSG = ""
		return nil, err
	}

	err = builder.setNextVersion(tags[0].Name)
	if err != nil {
		builder.spinner.FinalMSG = ""
		return nil, err
	}

	builder.tags = append(builder.tags, tags...)

	c := NewChangelog(builder.github.GetRepoName(), builder.github.GetRepoOwner())

	err = builder.buildChangeLog(c)
	if err != nil {
		builder.spinner.FinalMSG = ""
		return nil, err
	}

	return c, nil
}

func (builder *changelogBuilder) buildChangeLog(changelog Changelog) error {
	if configuration.Config.ShowUnreleased && builder.nextVersion == "" {
		builder.spinner.Suffix = " Getting unreleased entries"

		nextTag := builder.tags[0]
		pullRequests, err := builder.github.GetPullRequestsBetweenDates(nextTag.Date, Now())
		if err != nil {
			return err
		}

		unreleased := builder.populateUnreleasedEntry(pullRequests)
		if err != nil {
			return fmt.Errorf("could not process pull requests: %v", err)
		}

		changelog.AddUnreleased(unreleased)
	}

	for idx, currentTag := range builder.tags {
		builder.spinner.Suffix = fmt.Sprintf(" Processing tags: 🏷️  %s", currentTag.Name)
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

		entry, err := builder.populateReleasedEntry(
			currentTag.Name,
			previousTag.Name,

			currentTag.Date,
			pullRequests,
		)
		if err != nil {
			return fmt.Errorf("could not process pull requests: %v", err)
		}

		//changelog.entries = append(changelog.entries, *entry)
		changelog.AddEntry(*entry)
	}

	return nil
}

func (builder *changelogBuilder) populateUnreleasedEntry(pullRequests []githubclient.PullRequest) []string {
	unreleased := []string{}
	excludedLabels := configuration.Config.ExcludedLabels
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

			unreleased = append(unreleased, line)
		}
	}

	return unreleased
}

func (builder *changelogBuilder) populateReleasedEntry(currentTag string, previousTag string, date time.Time, pullRequests []githubclient.PullRequest) (*Entry, error) {
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

	excludedLabels := configuration.Config.ExcludedLabels
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
				err := entry.Append(section, line)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return entry, nil
}

func (builder *changelogBuilder) setNextVersion(currentVersion string) error {
	if builder.nextVersion != "" {
		if !utils.IsValidSemanticVersion(builder.nextVersion) {
			return fmt.Errorf("'%s' is not a valid semantic version", builder.nextVersion)
		}

		if !utils.VersionIsGreaterThan(currentVersion, builder.nextVersion) {
			return fmt.Errorf("the next version should be greater than the former: '%s' ≤ '%s'", builder.nextVersion, currentVersion)
		}

		lastCommitSha, err := builder.git.GetLastCommit()
		if err != nil {
			return err
		}

		tag := githubclient.Tag{
			Name: builder.nextVersion,
			Sha:  lastCommitSha,
			Date: Now(),
		}

		builder.tags = append(builder.tags, tag)
	}

	return nil
}

func hasExcludedLabel(excludedLabels []string, pr githubclient.PullRequest) bool {
	for _, label := range pr.Labels {
		if utils.SliceContainsString(excludedLabels, label.Name) {
			return true
		}
	}

	return false
}

func getSection(labels []githubclient.PullRequestLabel) string {
	sections := configuration.Config.Sections

	// Refactor
	lookup := make(map[string]string)
	for k, v := range sections {
		for _, label := range v {
			lookup[label] = k
		}
	}

	var section string
	skipUnlabelledEntries := configuration.Config.SkipEntriesWithoutLabel

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
