package changelog

import (
	"fmt"
	"sort"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chelnak/gh-changelog/internal/pkg/githubclient"
	"github.com/chelnak/gh-changelog/internal/pkg/utils"
	"github.com/google/go-github/v43/github"
	"github.com/spf13/viper"
)

func MakeFullChangelog(spinner *spinner.Spinner) (*ChangeLogProperties, error) {
	client, err := githubclient.NewGitHubClient()
	if err != nil {
		return nil, err
	}

	changeLog := NewChangeLogProperties(client.RepoContext.Owner, client.RepoContext.Repo)

	spinner.Suffix = " Gathering all tags"
	spinner.Start()

	tags, err := client.GetTags()
	if err != nil {
		return nil, err
	}

	// Sort by date, this is mad slow but you can't get at the date
	// any other way or sort the api response.
	spinner.Suffix = " Sorting entries"
	sort.Slice(tags, func(i, j int) bool {
		tagICommit, _ := client.GetCommit(tags[i].GetCommit().GetSHA())
		tagIDate := tagICommit.GetCommitter().GetDate()

		tagJCommit, _ := client.GetCommit(tags[j].GetCommit().GetSHA())
		tagJDate := tagJCommit.GetCommitter().GetDate()

		return tagIDate.After(tagJDate)
	})

	spinner.Suffix = " Gathering all pull requests"
	for idx, tag := range tags {
		spinner.Suffix = fmt.Sprintf(" Processing tags: 🏷️  %s", tag.GetName())
		currentCommit, err := client.GetCommit(tag.GetCommit().GetSHA())
		if err != nil {
			return nil, err
		}

		var nextCommit *github.Commit
		if idx+1 == len(tags) {
			nextCommit, err = client.GetFirstCommit()
		} else {
			nextCommit, err = client.GetCommit(tags[idx+1].GetCommit().GetSHA())
		}

		if err != nil {
			return nil, err
		}

		pullRequests, err := client.GetPullRequestsBetweenDates(
			nextCommit.GetCommitter().GetDate(),
			currentCommit.GetCommitter().GetDate(),
		)
		if err != nil {
			return nil, err
		}

		typeMap, err := processPullRequests(
			tag.GetName(),
			currentCommit.GetCommitter().GetDate(),
			pullRequests,
			viper.GetStringSlice("excludedLabels"),
		)
		if err != nil {
			return nil, err
		}

		changeLog.Tags = append(changeLog.Tags, *typeMap)
	}

	return changeLog, nil
}

func processPullRequests(tagName string, date time.Time, pullRequests []*github.Issue, excludedLabels []string) (*TagProperties, error) {
	// How to do this better?
	// This whole method is pretty grim.

	sections := viper.GetStringMapStringSlice("sections")

	lookup := make(map[string]string)
	for k, v := range sections {
		for _, label := range v {
			lookup[label] = k
		}
	}

	typeMap := NewTagProperties(tagName, date)
	for _, pr := range pullRequests {
		// Removed for now. This should be configurable in the future,
		// if len(pr.Labels) == 0 {
		// 	return nil, fmt.Errorf("could not process Pull Request #%d. All Pull Requests must have a label", pr.GetNumber())
		// }

		if !hasExcludedLabel(excludedLabels, pr) {
			entry := fmt.Sprintf(
				"%s [#%d](https://github.com/%s/%s/pull/%d) ([%s](https://github.com/%s))\n",
				pr.GetTitle(),
				pr.GetNumber(),
				pr.GetRepository().GetOwner().GetLogin(),
				pr.GetRepository().GetName(),
				pr.GetNumber(),
				pr.GetUser().GetLogin(),
				pr.GetUser().GetLogin(),
			)

			for _, label := range pr.Labels {
				if _, ok := lookup[label.GetName()]; ok {
					err := typeMap.Append(lookup[label.GetName()], entry)
					if err != nil {
						return nil, err
					}
				} else {
					// Add tp the "Other" section
					err := typeMap.Append("Other", entry)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return typeMap, nil
}

func hasExcludedLabel(excludedLabels []string, pr *github.Issue) bool {
	for _, label := range pr.Labels {
		if utils.Contains(excludedLabels, label.GetName()) {
			return true
		}
	}

	return false
}
