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
		return nil, fmt.Errorf("❌ %s", err)
	}

	changeLog := NewChangeLogProperties(client.RepoContext.Owner, client.RepoContext.Name)

	spinner.Suffix = " Gathering all tags"
	spinner.Start()

	tags, err := client.GetTags()
	if err != nil {
		return nil, fmt.Errorf("❌ could not get tags: %v", err)
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
			return nil, fmt.Errorf("❌ could not get commit for tag '%s': %v", tag.GetName(), err)
		}

		var nextCommit *github.Commit
		var nextTag string
		if idx+1 == len(tags) {
			nextCommit, err = client.GetFirstCommit()
			nextTag = nextCommit.GetSHA()
		} else {
			nextTag = tags[idx+1].GetName()
			nextCommit, err = client.GetCommit(tags[idx+1].GetCommit().GetSHA())
		}

		if err != nil {
			return nil, fmt.Errorf("❌ could not get next commit: %v", err)
		}

		pullRequests, err := client.GetPullRequestsBetweenDates(
			nextCommit.GetCommitter().GetDate(),
			currentCommit.GetCommitter().GetDate(),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"❌ could not get pull requests for range '%s - %s': %v",
				nextCommit.GetCommitter().GetDate(),
				currentCommit.GetCommitter().GetDate(),
				err,
			)
		}

		tagProperties, err := getTagProperties(
			tag.GetName(),
			nextTag,
			currentCommit.GetCommitter().GetDate(),
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
