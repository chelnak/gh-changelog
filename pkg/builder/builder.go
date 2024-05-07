// Package builder is responsible for building the changelog.
package builder

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/gitclient"
	"github.com/chelnak/gh-changelog/internal/githubclient"
	"github.com/chelnak/gh-changelog/internal/log"
	"github.com/chelnak/gh-changelog/internal/version"
	"github.com/chelnak/gh-changelog/pkg/changelog"
	"github.com/chelnak/gh-changelog/pkg/entry"
)

var Now = time.Now // must be a better way to stub this

type Git interface {
	GetFirstCommit() (string, error)
	GetLastCommit() (string, error)
	GetDateOfHash(hash string) (time.Time, error)
	GetCurrentBranch() (string, error)
	FetchAll() error
	Tags() ([]string, error)
	IsAncestorOf(commit, ancestor string) (bool, error)
}

type GitHub interface {
	GetPullRequestsBetweenDates(from, to time.Time) ([]githubclient.PullRequest, error)
	GetRepoName() string
	GetRepoOwner() string
}

type Builder struct {
	nextVersion   string
	fromVersion   string
	latestVersion bool
	latestTag     Tag
	tags          []Tag
	changelog     *changelog.Changelog
	git           Git
	github        GitHub
	branch        string
	ancestorsOnly bool
	filter        *regexp.Regexp
}

func NewBuilder() (*Builder, error) {
	builder := &Builder{}
	github, err := githubclient.NewGitHubClient()
	if err != nil {
		return nil, err
	}
	builder.github = github

	builder.git = gitclient.NewGitClient(exec.Command)
	builder.changelog = changelog.NewChangelog(github.GetRepoOwner(), github.GetRepoName())
	return builder, nil
}

func (b *Builder) NextVersion(nextVersion string) {
	b.nextVersion = nextVersion
}

func (b *Builder) FromVersion(fromVersion string) {
	b.fromVersion = fromVersion
}

func (b *Builder) LatestVersion() {
	b.latestVersion = true
}

func (b *Builder) Filter(filter *regexp.Regexp) {
	b.filter = filter
}

func (b *Builder) AncestorsOnly() {
	b.ancestorsOnly = true
}

func (b *Builder) Changelog(changelog *changelog.Changelog) {
	b.changelog = changelog
}

func (b *Builder) GitClient(client Git) {
	b.git = client
}

func (b *Builder) GitHubClient(client GitHub) {
	b.github = client
}

func (b *Builder) Build() (*changelog.Changelog, error) {
	defer log.Complete()
	log.Infof("Fetching tags...")
	if err := b.addTags(); err != nil {
		log.Errorf("could not fetch tags: %v", err)
		return nil, err
	}

	if b.nextVersion != "" {
		if err := b.setNextVersion(); err != nil {
			log.Errorf("could not set next version: %v", err)
			return nil, err
		}
	}

	if configuration.Config.ShowUnreleased && b.nextVersion == "" {
		log.Infof("Getting unreleased entries")
		err := b.getUnreleasedEntries()
		if err != nil {
			log.Errorf("could not process unreleased pull requests: %v", err)
			return nil, fmt.Errorf("could not process unreleased pull requests: %v", err)
		}
	}

	for i := 0; i < len(b.tags); i++ {
		var previousTag Tag
		if i+1 == len(b.tags) {
			previousTag = Tag{}
		} else {
			previousTag = b.tags[i+1]
		}

		err := b.getReleasedEntries(previousTag, b.tags[i])
		if err != nil {
			return nil, fmt.Errorf("could not process pull requests: %v", err)
		}
	}

	log.Infof("Open %s or run 'gh changelog show' to view your changelog.", configuration.Config.FileName)
	return b.changelog, nil
}

type Tag struct {
	Name string
	Sha  string
	Date time.Time
}

func (b *Builder) addTags() error {
	rawTags, err := b.git.Tags()
	if err != nil {
		return fmt.Errorf("updating tags: %w", err)
	}

	if len(rawTags) == 0 && b.nextVersion == "" {
		return errors.New("there are no tags on this repository to evaluate and the --next-version flag was not provided")
	}

	var tags []Tag
	for _, t := range rawTags {
		if t == "" {
			continue
		}

		parsedTag, err := b.parseTag(t)
		if err != nil {
			return fmt.Errorf("parsing tag %s: %w", t, err)
		}

		tags = append(tags, parsedTag)
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Date.After(tags[j].Date)
	})

	// Save the latest tag, so we can use it to get the unreleased entries
	b.latestTag = tags[0]

	if b.fromVersion != "" {
		tags = b.filterTagsFromVersion(tags)
	}

	if b.filter != nil {
		tags = b.filterTagsByName(tags, b.filter)
	}

	if b.ancestorsOnly {
		b.branch, err = b.git.GetCurrentBranch()
		if err != nil {
			return fmt.Errorf("fetching the current branch: %w", err)
		}

		tags, err = b.filterTagsForBranch(tags, b.branch)
		if err != nil {
			return fmt.Errorf("filtering tags for current branch: %w", err)
		}
	}

	if b.latestVersion {
		tags = []Tag{tags[0]}
	}

	b.tags = append(b.tags, tags...)
	return nil
}

func (b *Builder) parseTag(tag string) (Tag, error) {
	line := strings.Trim(tag, "'")
	line = strings.Trim(line, " ")
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		date, err := b.git.GetDateOfHash(parts[1])
		if err != nil {
			return Tag{}, fmt.Errorf("could not parse tag date from hash: %s: %v", parts[1], err)
		}

		return Tag{
			Name: parts[0],
			Sha:  parts[1],
			Date: date,
		}, nil
	}

	date, err := time.Parse(time.RFC3339, parts[2])
	if err != nil {
		return Tag{}, fmt.Errorf("could not parse tag date: %s: %v", parts[2], err)
	}

	return Tag{
		Name: parts[0],
		Sha:  parts[1],
		Date: date,
	}, nil
}

func (b *Builder) filterTagsFromVersion(tags []Tag) []Tag {
	for i, tag := range tags {
		if strings.EqualFold(tag.Name, b.fromVersion) {
			return tags[:i+1]
		}
	}
	return tags
}

func (b *Builder) filterTagsByName(tags []Tag, filter *regexp.Regexp) []Tag {
	var filteredTags []Tag
	for _, tag := range tags {
		if !filter.MatchString(tag.Name) {
			continue
		}
		filteredTags = append(filteredTags, tag)
	}
	return filteredTags
}

func (b *Builder) filterTagsForBranch(tags []Tag, branch string) ([]Tag, error) {
	var filteredTags []Tag
	for _, tag := range tags {
		if b.filter != nil {
			if b.filter.MatchString(tag.Name) {
				continue
			}
		}

		ancestor, err := b.git.IsAncestorOf(tag.Name, branch)
		if err != nil {
			log.Errorf("could not determine if %s is an ancestor of %s: %v", tag.Name, branch, err)
			return nil, err
		}
		if ancestor {
			filteredTags = append(filteredTags, tag)
		}
	}

	return filteredTags, nil
}

func (b *Builder) setNextVersion() error {
	if !version.IsValidSemanticVersion(b.nextVersion) {
		return fmt.Errorf("'%s' is not a valid semantic version", b.nextVersion)
	}
	if len(b.tags) > 0 {
		currentVersion := b.tags[0].Name
		if !version.NextVersionIsGreaterThanCurrent(b.nextVersion, currentVersion) {
			return fmt.Errorf("the next version should be greater than the former: '%s' ≤ '%s'", b.nextVersion, currentVersion)
		}
	}

	lastCommitSha, err := b.git.GetLastCommit()
	if err != nil {
		return fmt.Errorf("finding last commit: %w", err)
	}

	tag := Tag{
		Name: b.nextVersion,
		Sha:  lastCommitSha,
		Date: Now(),
	}

	b.tags = append([]Tag{tag}, b.tags...)

	return nil
}

func (b *Builder) getUnreleasedEntries() error {
	pullRequests, err := b.github.GetPullRequestsBetweenDates(b.latestTag.Date, Now())
	if err != nil {
		return err
	}

	var unreleased []string
	for _, pr := range pullRequests {
		if !hasExcludedLabel(pr) {
			line := b.formatEntryLine(pr)
			unreleased = append(unreleased, line)
		}
	}
	b.changelog.AddUnreleased(unreleased)
	return nil
}

func (b *Builder) getReleasedEntries(previousTag, currentTag Tag) error {
	log.Infof("Processing tag: 🏷️  %s", currentTag.Name)

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

func (b *Builder) formatEntryLine(pr githubclient.PullRequest) string {
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
		if slices.Contains(excludedLabels, label.Name) {
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
