package githubclient

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cli/go-gh"
	"github.com/shurcooL/githubv4"
)

type repoContext struct {
	owner string
	name  string
}

type GitHubClient interface {
	GetTags() ([]Tag, error)
	GetPullRequestsBetweenDates(from, to time.Time) ([]PullRequest, error)
	GetRepoName() string
	GetRepoOwner() string
}

type githubClient struct {
	base        *githubv4.Client
	repoContext repoContext
	httpContext context.Context
}

type Tag struct {
	Name string
	Sha  string
	Date time.Time
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

var TagQuery struct {
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

func (client *githubClient) GetTags() ([]Tag, error) {
	variables := map[string]interface{}{
		"repositoryOwner": githubv4.String(client.repoContext.owner),
		"repositoryName":  githubv4.String(client.repoContext.name),
		"cursor":          (*githubv4.String)(nil),
	}

	var tags []Tag

	for {
		err := client.base.Query(client.httpContext, &TagQuery, variables)
		if err != nil {
			return nil, fmt.Errorf("error getting tags: %w", err)
		}

		fmt.Println(TagQuery)

		for _, node := range TagQuery.Repository.Refs.Nodes {
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

		if !TagQuery.Repository.Refs.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = TagQuery.Repository.Refs.PageInfo.EndCursor
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Date.After(tags[j].Date)
	})

	return tags, nil
}

func (client *githubClient) GetPullRequestsBetweenDates(fromDate, toDate time.Time) ([]PullRequest, error) {
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
		"query":  githubv4.String(fmt.Sprintf(`repo:%s/%s is:pr is:merged merged:%s..%s`, client.repoContext.owner, client.repoContext.name, fromDate.Local().Format(time.RFC3339), toDate.Local().Format(time.RFC3339))),
		"cursor": (*githubv4.String)(nil),
	}

	var pullRequests []PullRequest

	for {
		err := client.base.Query(client.httpContext, &pullRequestSearchQuery, variables)
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

func (client *githubClient) GetRepoName() string {
	return client.repoContext.name
}

func (client *githubClient) GetRepoOwner() string {
	return client.repoContext.owner
}

func NewGitHubClient() (GitHubClient, error) {
	httpClient, err := gh.HTTPClient(nil)
	if err != nil {
		return nil, fmt.Errorf("could not create initial client: %s", err)
	}

	g := githubv4.NewClient(httpClient)

	currentRepository, err := gh.CurrentRepository()
	if err != nil {
		if strings.Contains(err.Error(), "not a git repository (or any of the parent directories)") {
			return nil, fmt.Errorf("the current directory is not a git repository")
		}

		return nil, err
	}

	client := &githubClient{
		base: g,
		repoContext: repoContext{
			owner: currentRepository.Owner(),
			name:  currentRepository.Name(),
		},
		httpContext: context.Background(),
	}

	return client, nil
}
