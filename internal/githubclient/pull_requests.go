package githubclient

import (
	"fmt"
	"time"

	graphql "github.com/cli/shurcooL-graphql"
)

type PullRequestLabel struct {
	Name string
}

type PullRequestEdge struct {
	Node struct {
		PullRequest struct {
			Number int
			Title  string
			Author struct {
				Login string
			}
			Labels struct {
				Nodes []PullRequestLabel
			} `graphql:"labels(first: 100)"`
		} `graphql:"... on PullRequest"`
	}
}

type PullRequestSearchQuery struct {
	Search struct {
		Edges    []PullRequestEdge
		PageInfo struct {
			EndCursor   graphql.String
			HasNextPage bool
		}
	} `graphql:"search(query: $query, type: ISSUE, first: 100, after: $cursor)"`
}

type PullRequest struct {
	Number int
	Title  string
	User   string
	Labels []PullRequestLabel
}

func (client *GitHub) GetPullRequestsBetweenDates(fromDate, toDate time.Time) ([]PullRequest, error) {
	variables := map[string]interface{}{
		"query": graphql.String(
			fmt.Sprintf(
				`repo:%s/%s is:pr is:merged merged:%s..%s`,
				client.repoContext.owner,
				client.repoContext.name,
				fromDate.Local().Format(time.RFC3339),
				toDate.Local().Format(time.RFC3339),
			),
		),
		"cursor": (*graphql.String)(nil),
	}

	var pullRequestSearchQuery PullRequestSearchQuery
	var pullRequests []PullRequest
	var edges []PullRequestEdge

	for {
		err := client.base.Query("pr", &pullRequestSearchQuery, variables)
		if err != nil {
			return nil, err
		}

		edges = append(edges, pullRequestSearchQuery.Search.Edges...)

		if !pullRequestSearchQuery.Search.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = pullRequestSearchQuery.Search.PageInfo.EndCursor
	}

	for _, edge := range edges {
		pullRequests = append(pullRequests, PullRequest{
			Number: edge.Node.PullRequest.Number,
			Title:  edge.Node.PullRequest.Title,
			User:   edge.Node.PullRequest.Author.Login,
			Labels: edge.Node.PullRequest.Labels.Nodes,
		})
	}

	return pullRequests, nil
}
