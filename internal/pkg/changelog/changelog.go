package changelog

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chelnak/gh-changelog/internal/pkg/gitclient"
	"github.com/chelnak/gh-changelog/internal/pkg/githubclient"
	"github.com/chelnak/gh-changelog/internal/pkg/utils"
	"github.com/google/go-github/v43/github"
	"github.com/spf13/viper"
)

func MakeFullChangelog(spinner *spinner.Spinner) (*ChangeLogProperties, error) {
	client, err := githubclient.NewGitHubClient()
	if err != nil {
		return nil, fmt.Errorf("❌ %s", err)
	}

	gitClient, err := gitclient.NewGitClient()
	if err != nil {
		return nil, fmt.Errorf("❌ %s", err)
	}

	changeLog := NewChangeLogProperties(client.RepoContext.Owner, client.RepoContext.Name)

	spinner.Suffix = " Gathering all tags"
	spinner.Start()

	tags, err := gitClient.GetTags()
	if err != nil {
		return nil, fmt.Errorf("❌ could not get tags: %v", err)
	}

	spinner.Suffix = " Gathering all pull requests"
	for idx, currentTag := range tags {
		spinner.Suffix = fmt.Sprintf(" Processing tags: 🏷️  %s", currentTag.Name)

		var nextTag *gitclient.Ref
		if idx+1 == len(tags) {
			nextTag, err = gitClient.GetFirstCommit()
			if err != nil {
				return nil, fmt.Errorf("❌ could not get first commit: %v", err)
			}
		} else {
			nextTag = tags[idx+1]
		}

		pullRequests, err := client.GetPullRequestsBetweenDates(nextTag.Date, currentTag.Date)
		if err != nil {
			return nil, fmt.Errorf(
				"❌ could not get pull requests for range '%s - %s': %v",
				nextTag.Date,
				currentTag.Date,
				err,
			)
		}

		tagProperties, err := getTagProperties(
			currentTag.Name,
			nextTag.Name,
			currentTag.Date,
			pullRequests,
			viper.GetStringSlice("excludedLabels"),
			client.RepoContext,
		)
		if err != nil {
			return nil, fmt.Errorf("❌ could not process pull requests: %v", err)
		}

		changeLog.Tags = append(changeLog.Tags, *tagProperties)
	}

	return changeLog, nil
}

func getTagProperties(currentTag string, nextTag string, date time.Time, pullRequests []*github.Issue, excludedLabels []string, repoContext githubclient.RepoContext) (*TagProperties, error) {
	tagProperties := NewTagProperties(currentTag, nextTag, date)
	for _, pr := range pullRequests {
		if !hasExcludedLabel(excludedLabels, pr) {
			entry := fmt.Sprintf(
				"%s [#%d](https://github.com/%s/%s/pull/%d) ([%s](https://github.com/%s))\n",
				pr.GetTitle(),
				pr.GetNumber(),
				repoContext.Owner,
				repoContext.Name,
				pr.GetNumber(),
				pr.GetUser().GetLogin(),
				pr.GetUser().GetLogin(),
			)

			section := getSection(pr.Labels)
			err := tagProperties.Append(section, entry)
			if err != nil {
				return nil, fmt.Errorf("❌ could not append entry: %v", err)
			}
		}
	}

	return tagProperties, nil
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

	section := "Other"
	for _, label := range labels {
		if _, ok := lookup[label.GetName()]; ok {
			section = lookup[label.GetName()]
		}
	}

	return section
}
