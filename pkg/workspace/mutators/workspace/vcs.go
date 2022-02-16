package workspace

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/mify-io/mify/pkg/workspace/mutators"
)

func getCommitOpts(repo *git.Repository) (git.CommitOptions, error) {
	o := git.CommitOptions{}
	cfg, err := repo.ConfigScoped(config.SystemScope)
	if err != nil {
		return o, err
	}

	if cfg.Author.Email != "" && cfg.Author.Name != "" {
		o.Author = &object.Signature{
			Name:  cfg.Author.Name,
			Email: cfg.Author.Email,
			When:  time.Now(),
		}
	}

	if cfg.Committer.Email != "" && cfg.Committer.Name != "" {
		o.Committer = &object.Signature{
			Name:  cfg.Committer.Name,
			Email: cfg.Committer.Email,
			When:  time.Now(),
		}
	}

	if o.Author == nil && cfg.User.Email != "" && cfg.User.Name != "" {
		o.Author = &object.Signature{
			Name:  cfg.User.Name,
			Email: cfg.User.Email,
			When:  time.Now(),
		}
	}
	if o.Author == nil {
		o.Author = &object.Signature{
			Name:  "Mify",
			Email: "support@mify.io",
			When:  time.Now(),
		}
	}
	return o, nil
}

func InitGit(mutContext *mutators.MutatorContext) error {
	repo, err := git.PlainInit(mutContext.GetDescription().BasePath, false)
	if err != nil {
		return err
	}
	tree, err := repo.Worktree()
	if err != nil {
		return err
	}

	opts, err := getCommitOpts(repo)
	if err != nil {
		return err
	}

	_, err = tree.Commit("Initial commit", &opts)
	return err
}
