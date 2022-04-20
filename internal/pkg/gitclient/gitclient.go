package gitclient

import (
	"sort"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Ref struct {
	Sha  string
	Name string
	Date time.Time
}

type GitClient struct {
	repo *git.Repository
}

func (gc *GitClient) GetTags() ([]*Ref, error) {
	refs, err := gc.repo.References()
	if err != nil {
		return nil, err
	}

	tags := []*Ref{}
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsTag() {
			hash, err := gc.repo.ResolveRevision(plumbing.Revision(ref.Hash().String()))
			if err != nil {
				return err
			}

			commit, err := gc.repo.CommitObject(*hash)
			if err != nil {
				return err
			}

			tags = append(tags, &Ref{
				Sha:  hash.String(),
				Name: ref.Name().Short(),
				Date: commit.Committer.When,
			})

			if err != nil {
				return err
			}
		}

		return nil
	})

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Date.After(tags[j].Date)
	})

	return tags, err
}

func (gc *GitClient) GetFirstCommit() (*Ref, error) {
	commitObjectss, err := gc.repo.CommitObjects()
	if err != nil {
		return nil, err
	}

	commits := []*object.Commit{}
	err = commitObjectss.ForEach(func(commit *object.Commit) error {
		commits = append(commits, commit)
		return nil
	})

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Committer.When.Before(commits[j].Committer.When)
	})

	if err != nil {
		return nil, err
	}

	ref := &Ref{
		Sha:  commits[0].Hash.String(),
		Name: commits[0].Hash.String(), // Set Name to the commit hash because it is the first commit and doesn't have a name
		Date: commits[0].Committer.When,
	}

	return ref, nil
}

func NewGitClient() (*GitClient, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		panic(err)
	}

	// TODO: Make this configurable
	err = repo.Fetch(&git.FetchOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, err
	}

	return &GitClient{
		repo: repo,
	}, nil
}
