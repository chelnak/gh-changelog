// Package builder is responsible for building the changelog.
package builder

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/gitclient"
	"github.com/chelnak/gh-changelog/internal/githubclient"
	"github.com/chelnak/gh-changelog/internal/logging"
	"github.com/chelnak/gh-changelog/internal/utils"
	"github.com/chelnak/gh-changelog/pkg/changelog"
	"github.com/chelnak/gh-changelog/pkg/entry"
)

var Now = time.Now // must be a better way to stub this

type BuilderOptions struct {
	Logger        string
	NextVersion   string
	FromVersion   string
	LatestVersion bool
	GitClient     gitclient.GitClient
	GitHubClient  githubclient.GitHubClient
}

func (bo *BuilderOptions) setupGitClient() {
	if bo.GitClient == nil {
		bo.GitClient = gitclient.NewGitClient(exec.Command)
	}
}

func (bo *BuilderOptions) setupGitHubClient() error {
	if bo.GitHubClient == nil {
		client, err := githubclient.NewGitHubClient()
		if err != nil {
			return err
		}
		bo.GitHubClient = client
	}

	return nil
}

type Builder interface {
	BuildChangelog() (changelog.Changelog, error)
}

type builder struct {
	nextVersion   string
	fromVersion   string
	latestVersion bool
	tags          []githubclient.Tag
	changelog     changelog.Changelog
	git           gitclient.GitClient
	github        githubclient.GitHubClient
	logger        logging.Logger
}

func NewBuilder(options BuilderOptions) (Builder, error) {
	options.setupGitClient()

	if err := options.setupGitHubClient(); err != nil {
		return nil, err
	}

	changelog := changelog.NewChangelog(
		options.GitHubClient.GetRepoOwner(),
		options.GitHubClient.GetRepoName(),
	)

	builder := &builder{
		nextVersion:   options.NextVersion,
		fromVersion:   options.FromVersion,
		latestVersion: options.LatestVersion,
		changelog:     changelog,
		git:           options.GitClient,
		github:        options.GitHubClient,
	}

	loggerType, err := logging.GetLoggerType(options.Logger)
	if err != nil {
		return builder, err
	}

	builder.logger = logging.NewLogger(loggerType)

	return builder, nil
}

func (b *builder) BuildChangelog() (changelog.Changelog, error) {
	// defer b.spinnerManager.Stop()

	b.logger.Infof("Fetching tags...")
	err := b.updateTags()
	if err != nil {
		b.logger.Errorf(err.Error())
		return nil, err
	}

	if b.nextVersion != "" {
		err = b.setNextVersion()
		if err != nil {
			b.logger.Errorf(err.Error())
			return nil, err
		}
	}

	if configuration.Config.ShowUnreleased && b.nextVersion == "" {
		b.logger.Infof("Getting unreleased entries")
		err := b.getUnreleasedEntries()
		if err != nil {
			return nil, err
		}
	}

	for i := 0; i < len(b.tags); i++ {
		var previousTag githubclient.Tag
		if i+1 == len(b.tags) {
			previousTag = githubclient.Tag{}
		} else {
			previousTag = b.tags[i+1]
		}

		err := b.getReleasedEntries(previousTag, b.tags[i])
		if err != nil {
			return nil, fmt.Errorf("could not process pull requests: %v", err)
		}

		if strings.EqualFold(b.fromVersion, b.tags[i].Name) || b.latestVersion {
			break
		}
	}

	b.logger.Infof("Open %s or run 'gh changelog show' to view your changelog.", configuration.Config.FileName)
	b.logger.Complete()

	return b.changelog, nil
}

func (b *builder) updateTags() error {
	tags, err := b.github.GetTags()
	if err != nil {
		return err
	}

	if len(tags) == 0 && b.nextVersion == "" {
		return errors.New("there are no tags on this repository to evaluate and the --next-version flag was not provided")
	}

	b.tags = append(b.tags, tags...)

	return nil
}

func (b *builder) setNextVersion() error {
	if !utils.IsValidSemanticVersion(b.nextVersion) {
		return fmt.Errorf("'%s' is not a valid semantic version", b.nextVersion)
	}
	if len(b.tags) > 0 {
		currentVersion := b.tags[0].Name
		if !utils.NextVersionIsGreaterThanCurrent(b.nextVersion, currentVersion) {
			return fmt.Errorf("the next version should be greater than the former: '%s' ≤ '%s'", b.nextVersion, currentVersion)
		}
	}

	lastCommitSha, err := b.git.GetLastCommit()
	if err != nil {
		return err
	}

	tag := githubclient.Tag{
		Name: b.nextVersion,
		Sha:  lastCommitSha,
		Date: Now(),
	}

	b.tags = append([]githubclient.Tag{tag}, b.tags...)

	return nil
}

func (b *builder) getUnreleasedEntries() error {
	pullRequests, err := b.github.GetPullRequestsBetweenDates(b.tags[0].Date, Now())
	if err != nil {
		return err
	}

	unreleased := []string{}
	for _, pr := range pullRequests {
		if !hasExcludedLabel(pr) {
			line := b.formatEntryLine(pr)
			unreleased = append(unreleased, line)
		}
	}

	b.changelog.AddUnreleased(unreleased)

	return nil
}

func (b *builder) getReleasedEntries(previousTag, currentTag githubclient.Tag) error {
	b.logger.Infof("Processing tag: 🏷️  %s", currentTag.Name)

	pullRequests, err := b.github.GetPullRequestsBetweenDates(previousTag.Date, currentTag.Date)
	if err != nil {
		return err
	}

	e := entry.NewEntry(currentTag.Name, currentTag.Date)

	for _, pr := range pullRequests {
		if !hasExcludedLabel(pr) {
			section := getSection(pr.Labels)
			line := b.formatEntryLine(pr)

			if section != "" {
				err := e.Append(section, line)
				if err != nil {
					return err
				}
			}
		}
	}
	b.changelog.Insert(e)
	return nil
}

func (b *builder) formatEntryLine(pr githubclient.PullRequest) string {
	return fmt.Sprintf(
		"%s [#%d](https://github.com/%s/%s/pull/%d) ([%s](https://github.com/%s))",
		pr.Title,
		pr.Number,
		b.github.GetRepoOwner(),
		b.github.GetRepoName(),
		pr.Number,
		pr.User,
		pr.User,
	)
}

func hasExcludedLabel(pr githubclient.PullRequest) bool {
	excludedLabels := configuration.Config.ExcludedLabels
	for _, label := range pr.Labels {
		if utils.SliceContainsString(excludedLabels, label.Name) {
			return true
		}
	}

	return false
}

func getSection(labels []githubclient.PullRequestLabel) string {
	sections := configuration.Config.Sections

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
