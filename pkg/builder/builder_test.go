package builder_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/githubclient"
	"github.com/chelnak/gh-changelog/internal/log"
	mocks "github.com/chelnak/gh-changelog/mocks/builder"
	"github.com/chelnak/gh-changelog/pkg/builder"
	"github.com/chelnak/gh-changelog/pkg/changelog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	repoName  = "repo-name"
	repoOwner = "repo-owner"
)

var rawTags = `v0.1.0 577eab3665f6c5c0220514c465fd7e7fff2d3d4f 2022-04-15T09:56:13+01:00
v0.2.0 aa101426c006214dd1df27ab4a545b0309d71d7c 2022-04-15T22:46:07+01:00`

func TestBuilder(t *testing.T) {
	log.SetupLogging(log.NOOP)
	_ = configuration.InitConfig()

	t.Run("should build a changelog", func(t *testing.T) {
		_ = configuration.InitConfig()

		b, err := builder.NewBuilder()
		require.NoError(t, err)

		mockGit := mocks.NewMockGit(t)
		mockGit.EXPECT().Tags().Return(strings.Split(rawTags, "\n"), nil)
		b.GitClient(mockGit)

		mockGitHub := mocks.NewMockGitHub(t)
		mockGitHub.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{}, nil).Once()
		mockGitHub.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{
			{
				Number: 2,
				Title:  "this is a test pr 2",
				User:   "test-user",
				Labels: []githubclient.PullRequestLabel{
					{
						Name: "enhancement",
					},
				},
			},
			{
				Number: 1,
				Title:  "this is a test pr",
				User:   "test-user",
				Labels: []githubclient.PullRequestLabel{
					{
						Name: "enhancement",
					},
				},
			},
		}, nil).Once()
		mockGitHub.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{}, nil).Once()

		mockGitHub.EXPECT().GetRepoName().Return(repoName)
		mockGitHub.EXPECT().GetRepoOwner().Return(repoOwner)
		b.GitHubClient(mockGitHub)

		b.Changelog(changelog.NewChangelog(repoOwner, repoName))

		c, err := b.Build()
		require.NoError(t, err)
		require.Equal(t, repoName, c.GetRepoName())
		require.Equal(t, repoOwner, c.GetRepoOwner())

		require.Len(t, c.GetUnreleased(), 0)
		require.Len(t, c.GetEntries(), 2)

		require.Equal(
			t,
			"this is a test pr 2 [#2](https://github.com/repo-owner/repo-name/pull/2) ([test-user](https://github.com/test-user))",
			c.GetEntries()[0].Added[0],
		)
	})

	t.Run("should error with an older next version", func(t *testing.T) {
		b, err := builder.NewBuilder()
		require.NoError(t, err)

		mockGit := mocks.NewMockGit(t)
		mockGit.EXPECT().Tags().Return(strings.Split(rawTags, "\n"), nil)

		b.GitClient(mockGit)

		mockGitHb := mocks.NewMockGitHub(t)
		b.GitHubClient(mockGitHb)

		b.Changelog(changelog.NewChangelog(repoOwner, repoName))

		b.NextVersion("v0.0.1")

		_, err = b.Build()
		require.Error(t, err)
		require.Equal(t, "the next version should be greater than the former: 'v0.0.1' ≤ 'v0.2.0'", err.Error())
	})

	t.Run("should error with no tags", func(t *testing.T) {
		b, err := builder.NewBuilder()
		require.NoError(t, err)

		mockGit := mocks.NewMockGit(t)
		mockGit.EXPECT().Tags().Return([]string{}, nil)
		b.GitClient(mockGit)

		mockGitHb := mocks.NewMockGitHub(t)
		b.GitHubClient(mockGitHb)

		b.Changelog(changelog.NewChangelog(repoOwner, repoName))

		_, err = b.Build()
		require.Error(t, err)
		require.Equal(t, "there are no tags on this repository to evaluate and the --next-version flag was not provided", err.Error())
	})

	t.Run("should build a changelog from a given version", func(t *testing.T) {
		b, err := builder.NewBuilder()
		require.NoError(t, err)

		mockGit := mocks.NewMockGit(t)
		mockGit.EXPECT().Tags().Return(strings.Split(rawTags, "\n"), nil)
		b.GitClient(mockGit)

		mockGitHub := mocks.NewMockGitHub(t)
		mockGitHub.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{}, nil).Once()
		mockGitHub.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{
			{
				Number: 2,
				Title:  "this is a test pr 2",
				User:   "test-user",
				Labels: []githubclient.PullRequestLabel{
					{
						Name: "enhancement",
					},
				},
			},
			{
				Number: 1,
				Title:  "this is a test pr",
				User:   "test-user",
				Labels: []githubclient.PullRequestLabel{
					{
						Name: "enhancement",
					},
				},
			},
		}, nil).Once()
		mockGitHub.EXPECT().GetRepoName().Return(repoName).Maybe()
		mockGitHub.EXPECT().GetRepoOwner().Return(repoOwner).Maybe()
		b.GitHubClient(mockGitHub)

		b.Changelog(changelog.NewChangelog(repoOwner, repoName))

		b.FromVersion("v0.2.0")

		c, err := b.Build()
		require.NoError(t, err)
		require.Equal(t, repoName, c.GetRepoName())
		require.Equal(t, repoOwner, c.GetRepoOwner())

		require.Len(t, c.GetUnreleased(), 0)
		require.Len(t, c.GetEntries(), 1)

		require.Equal(
			t,
			"this is a test pr 2 [#2](https://github.com/repo-owner/repo-name/pull/2) ([test-user](https://github.com/test-user))",
			c.GetEntries()[0].Added[0],
		)
	})

	t.Run("should build a changelog that contains only the last version", func(t *testing.T) {
		b, err := builder.NewBuilder()
		require.NoError(t, err)

		mockGit := mocks.NewMockGit(t)
		mockGit.EXPECT().Tags().Return(strings.Split(rawTags, "\n"), nil)
		b.GitClient(mockGit)

		mockGitHub := mocks.NewMockGitHub(t)
		mockGitHub.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{}, nil).Once()
		mockGitHub.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{
			{
				Number: 2,
				Title:  "this is a test pr 2",
				User:   "test-user",
				Labels: []githubclient.PullRequestLabel{
					{
						Name: "enhancement",
					},
				},
			},
			{
				Number: 1,
				Title:  "this is a test pr",
				User:   "test-user",
				Labels: []githubclient.PullRequestLabel{
					{
						Name: "enhancement",
					},
				},
			},
		}, nil)
		mockGitHub.EXPECT().GetRepoName().Return(repoName).Maybe()
		mockGitHub.EXPECT().GetRepoOwner().Return(repoOwner).Maybe()
		b.GitHubClient(mockGitHub)

		b.Changelog(changelog.NewChangelog(repoOwner, repoName))

		b.LatestVersion()

		c, err := b.Build()
		require.NoError(t, err)
		require.Equal(t, repoName, c.GetRepoName())
		require.Equal(t, repoOwner, c.GetRepoOwner())

		require.Len(t, c.GetUnreleased(), 0)
		require.Len(t, c.GetEntries(), 1)

		require.Equal(
			t,
			"this is a test pr 2 [#2](https://github.com/repo-owner/repo-name/pull/2) ([test-user](https://github.com/test-user))",
			c.GetEntries()[0].Added[0],
		)
	})

	t.Run("should build a changelog that contains only the last version and the next version", func(t *testing.T) {
		b, err := builder.NewBuilder()
		require.NoError(t, err)

		mockGit := mocks.NewMockGit(t)
		mockGit.EXPECT().Tags().Return(strings.Split(rawTags, "\n"), nil)
		mockGit.EXPECT().GetLastCommit().Return("aa101426c006214dd1df27ab4a545b0309d71d7c", nil)
		b.GitClient(mockGit)

		mockGitHub := mocks.NewMockGitHub(t)
		// mockGitHub.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{}, nil).Once()
		mockGitHub.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{
			{
				Number: 4,
				Title:  "this is a test pr 4",
				User:   "test-user",
				Labels: []githubclient.PullRequestLabel{
					{
						Name: "enhancement",
					},
				},
			},
			{
				Number: 3,
				Title:  "this is a test pr 3",
				User:   "test-user",
				Labels: []githubclient.PullRequestLabel{
					{
						Name: "enhancement",
					},
				},
			},
		}, nil).Once()
		mockGitHub.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{
			{
				Number: 2,
				Title:  "this is a test pr 2",
				User:   "test-user",
				Labels: []githubclient.PullRequestLabel{
					{
						Name: "enhancement",
					},
				},
			},
			{
				Number: 1,
				Title:  "this is a test pr",
				User:   "test-user",
				Labels: []githubclient.PullRequestLabel{
					{
						Name: "enhancement",
					},
				},
			},
		}, nil).Once()

		mockGitHub.EXPECT().GetRepoName().Return(repoName).Maybe()
		mockGitHub.EXPECT().GetRepoOwner().Return(repoOwner).Maybe()
		b.GitHubClient(mockGitHub)

		b.Changelog(changelog.NewChangelog(repoOwner, repoName))

		b.LatestVersion()
		b.NextVersion("v0.3.0")

		c, err := b.Build()
		require.NoError(t, err)
		require.Equal(t, repoName, c.GetRepoName())
		require.Equal(t, repoOwner, c.GetRepoOwner())

		require.Len(t, c.GetUnreleased(), 0)
		require.Len(t, c.GetEntries(), 2)

		require.Equal(
			t,
			"this is a test pr 4 [#4](https://github.com/repo-owner/repo-name/pull/4) ([test-user](https://github.com/test-user))",
			c.GetEntries()[0].Added[0],
		)

		require.Equal(
			t,
			"this is a test pr 2 [#2](https://github.com/repo-owner/repo-name/pull/2) ([test-user](https://github.com/test-user))",
			c.GetEntries()[1].Added[0],
		)
	})

	t.Run("should build a changelog that contains only the entries matching a given filter", func(t *testing.T) {
		var filterTags = `v0.1.0 577eab3665f6c5c0220514c465fd7e7fff2d3d4f 2022-04-15T09:56:13+01:00
v0.2.0 aa101426c006214dd1df27ab4a545b0309d71d7c 2022-04-16T22:46:07+01:00
v0.3.0 577eab3665f6c5c0220514c465fd7e7fff2d3d4f 2022-04-17T09:56:13+01:00
v0.4.0 aa101426c006214dd1df27ab4a545b0309d71d7c 2022-04-18T22:46:07+01:00`

		tests := []struct {
			name             string
			filter           *regexp.Regexp
			mockGitHubClient func() *mocks.MockGitHub
			expected         string
		}{
			{
				name:   "filter: .0.4.*",
				filter: regexp.MustCompile(".0.4.*"),
				mockGitHubClient: func() *mocks.MockGitHub {
					m := mocks.NewMockGitHub(t)
					m.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{}, nil).Once()
					m.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{
						{
							Number: 4,
							Title:  "this is a test pr 4",
							User:   "test-user",
							Labels: []githubclient.PullRequestLabel{
								{
									Name: "enhancement",
								},
							},
						},
					}, nil).Once()
					m.EXPECT().GetRepoName().Return(repoName)
					m.EXPECT().GetRepoOwner().Return(repoOwner)
					return m
				},
				expected: "this is a test pr 4 [#4](https://github.com/repo-owner/repo-name/pull/4) ([test-user](https://github.com/test-user))",
			},
			{
				name:   "filter: v0.4.*",
				filter: regexp.MustCompile("v0.4.*"),
				mockGitHubClient: func() *mocks.MockGitHub {
					m := mocks.NewMockGitHub(t)
					m.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{}, nil).Once()
					m.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{
						{
							Number: 4,
							Title:  "this is a test pr 4",
							User:   "test-user",
							Labels: []githubclient.PullRequestLabel{
								{
									Name: "enhancement",
								},
							},
						},
					}, nil).Once()
					m.EXPECT().GetRepoName().Return(repoName)
					m.EXPECT().GetRepoOwner().Return(repoOwner)
					return m
				},
				expected: "this is a test pr 4 [#4](https://github.com/repo-owner/repo-name/pull/4) ([test-user](https://github.com/test-user))",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				b, err := builder.NewBuilder()
				require.NoError(t, err)

				mockGit := mocks.NewMockGit(t)
				mockGit.EXPECT().Tags().Return(strings.Split(filterTags, "\n"), nil)
				b.GitClient(mockGit)

				b.GitHubClient(tt.mockGitHubClient())
				b.Changelog(changelog.NewChangelog(repoOwner, repoName))
				b.Filter(tt.filter)
				c, err := b.Build()
				require.NoError(t, err)

				require.Len(t, c.GetEntries(), 1)
				require.Equal(
					t,
					tt.expected,
					c.GetEntries()[0].Added[0],
				)
			})
		}
	})

	t.Run("should build a changelog that contains only ancestors of the current branch", func(t *testing.T) {
		var filterTags = `v0.1.0 577eab3665f6c5c0220514c465fd7e7fff2d3d4f 2022-04-15T09:56:13+01:00
v0.2.0 aa101426c006214dd1df27ab4a545b0309d71d7c 2022-04-16T22:46:07+01:00
v0.3.0 577eab3665f6c5c0220514c465fd7e7fff2d3d4f 2022-04-17T09:56:13+01:00
v0.4.0 aa101426c006214dd1df27ab4a545b0309d71d7c 2022-04-18T22:46:07+01:00`

		mockGitHubClient := mocks.NewMockGitHub(t)
		mockGitHubClient.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{}, nil).Once()
		mockGitHubClient.EXPECT().GetPullRequestsBetweenDates(mock.Anything, mock.Anything).Return([]githubclient.PullRequest{
			{
				Number: 4,
				Title:  "this is a test pr 4",
				User:   "test-user",
				Labels: []githubclient.PullRequestLabel{
					{
						Name: "enhancement",
					},
				},
			},
		}, nil).Once()
		mockGitHubClient.EXPECT().GetRepoName().Return(repoName)
		mockGitHubClient.EXPECT().GetRepoOwner().Return(repoOwner)

		mockGitClient := mocks.NewMockGit(t)
		mockGitClient.EXPECT().Tags().Return(strings.Split(filterTags, "\n"), nil)
		mockGitClient.EXPECT().GetCurrentBranch().Return("foo", nil)
		mockGitClient.EXPECT().IsAncestorOf(mock.Anything, mock.Anything).Return(false, nil).Once()
		mockGitClient.EXPECT().IsAncestorOf(mock.Anything, mock.Anything).Return(true, nil).Once()
		mockGitClient.EXPECT().IsAncestorOf(mock.Anything, mock.Anything).Return(false, nil).Once()
		mockGitClient.EXPECT().IsAncestorOf(mock.Anything, mock.Anything).Return(false, nil).Once()

		b, err := builder.NewBuilder()
		require.NoError(t, err)

		b.GitClient(mockGitClient)
		b.GitHubClient(mockGitHubClient)
		b.Changelog(changelog.NewChangelog(repoOwner, repoName))
		b.AncestorsOnly()
		c, err := b.Build()
		require.NoError(t, err)

		require.Len(t, c.GetEntries(), 1)
		require.Equal(
			t,
			"this is a test pr 4 [#4](https://github.com/repo-owner/repo-name/pull/4) ([test-user](https://github.com/test-user))",
			c.GetEntries()[0].Added[0],
		)
	})
}
