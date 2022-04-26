package githubv4client

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cli/go-gh"
	"github.com/shurcooL/githubv4"
)

type RepoContext struct {
	Owner string
	Name  string
}

type GitHubGraphClient struct {
	baseClient  *githubv4.Client
	RepoContext RepoContext
	httpContext context.Context
}

func NewGitHubGraphClient() (*GitHubGraphClient, error) {
	httpClient, err := gh.HTTPClient(nil)
	if err != nil {
		return nil, err
	}

	g := githubv4.NewClient(httpClient)

	currentRepository, err := gh.CurrentRepository()
	if err != nil {
		if strings.Contains(err.Error(), "not a git repository (or any of the parent directories)") {
			return nil, fmt.Errorf("the current directory is not a git repository")
		}

		return nil, err
	}

	client := &GitHubGraphClient{
		baseClient: g,
		RepoContext: RepoContext{
			Owner: currentRepository.Owner(),
			Name:  currentRepository.Name(),
		},
		httpContext: context.Background(),
	}

	return client, nil
}

type Tag struct {
	Name string
	Sha  string
	Date time.Time
}

func (c *GitHubGraphClient) GetTags() ([]Tag, error) {
	var tagQuery struct {
		Repository struct {
			Refs struct {
				Nodes []struct {
					Name   string
					Target struct {
						TypeName string `graphql:"__typename"`
						Tag      struct {
							Oid    string
							Tagger struct {
								Date time.Time
							}
						} `graphql:"... on Tag"`
						Commit struct {
							Oid       string
							Committer struct {
								Date time.Time
							}
						} `graphql:"... on Commit"`
					}
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"refs(refPrefix: \"refs/tags/\", last: 100, after: $cursor)"`
		} `graphql:"repository(owner:$repositoryOwner,name:$repositoryName)"`
	}

	variables := map[string]interface{}{
		"repositoryOwner": githubv4.String(c.RepoContext.Owner),
		"repositoryName":  githubv4.String(c.RepoContext.Name),
		"cursor":          (*githubv4.String)(nil),
	}

	var tags []Tag

	for {
		err := c.baseClient.Query(c.httpContext, &tagQuery, variables)
		if err != nil {
			return nil, err
		}

		for _, node := range tagQuery.Repository.Refs.Nodes {
			switch node.Target.TypeName {
			case "Tag":
				tags = append(tags, Tag{
					Name: node.Name,
					Sha:  node.Target.Tag.Oid,
					Date: node.Target.Tag.Tagger.Date,
				})
			case "Commit":
				tags = append(tags, Tag{
					Name: node.Name,
					Sha:  node.Target.Commit.Oid,
					Date: node.Target.Commit.Committer.Date,
				})
			}
		}

		if !tagQuery.Repository.Refs.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = tagQuery.Repository.Refs.PageInfo.EndCursor
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Date.After(tags[j].Date)
	})

	return tags, nil
}

type Label struct {
	Name string
}

type PullRequest struct {
	Number int
	Title  string
	User   string
	Labels []Label
}

func (c *GitHubGraphClient) GetPullRequestsBetweenDates(fromDate time.Time, toDate time.Time) ([]PullRequest, error) {
	var pullRequestSearchQuery struct {
		Search struct {
			Edges []struct {
				Node struct {
					PullRequest struct {
						Number int
						Title  string
						Author struct {
							Login string
						}
						Labels struct {
							Nodes []struct {
								Name string
							}
						} `graphql:" labels(first: 100)"`
					} `graphql:"... on PullRequest"`
				}
			}
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"search(query: $query, type: ISSUE, first: 100, after: $cursor)"`
	}

	variables := map[string]interface{}{
		"query":  githubv4.String(fmt.Sprintf(`repo:%s/%s is:pr is:merged merged:%s..%s`, c.RepoContext.Owner, c.RepoContext.Name, fromDate.Local().Format(time.RFC3339), toDate.Local().Format(time.RFC3339))),
		"cursor": (*githubv4.String)(nil),
	}

	var pullRequests []PullRequest

	for {
		err := c.baseClient.Query(c.httpContext, &pullRequestSearchQuery, variables)
		if err != nil {
			return nil, err
		}

		for _, edge := range pullRequestSearchQuery.Search.Edges {
			pullRequests = append(pullRequests, PullRequest{
				Number: edge.Node.PullRequest.Number,
				Title:  edge.Node.PullRequest.Title,
				User:   edge.Node.PullRequest.Author.Login,
				Labels: make([]Label, len(edge.Node.PullRequest.Labels.Nodes)),
			})

			for i, label := range edge.Node.PullRequest.Labels.Nodes {
				pullRequests[len(pullRequests)-1].Labels[i] = Label{
					Name: label.Name,
				}
			}
		}

		if !pullRequestSearchQuery.Search.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = pullRequestSearchQuery.Search.PageInfo.EndCursor
	}
	return pullRequests, nil
}
