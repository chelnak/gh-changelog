package changelog

import (
	"fmt"

	"github.com/chelnak/gh-changelog/internal/pkg/gitclient"
	"github.com/chelnak/gh-changelog/internal/pkg/githubclient"
)

func NewChangelog() (*ChangeLog, error) {
	githubClient, err := githubclient.NewGitHubClient()
	if err != nil {
		return nil, fmt.Errorf("❌ %s", err)
	}

	gitClient, err := gitclient.NewGitClient()
	if err != nil {
		return nil, fmt.Errorf("❌ %s", err)
	}

	tags, err := gitClient.GetTags()
	if err != nil {
		return nil, fmt.Errorf("❌ could not get tags: %v", err)
	}

	if len(tags) < 1 {
		return nil, fmt.Errorf("💡 no tags found in the current repository")
	}

	builder := NewChangeLogBuilder(gitClient, githubClient, tags)
	changeLog, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("❌ %s", err)
	}

	return changeLog, nil
}
