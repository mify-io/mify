package workspace

import (
	"github.com/go-git/go-git/v5"
	"github.com/mify-io/mify/pkg/workspace/mutators"
)

func InitGit(mutContext *mutators.MutatorContext) error {
	repo, err := git.PlainInit(mutContext.GetDescription().BasePath, false)
	if err != nil {
		return err
	}
	tree, err := repo.Worktree()
	if err != nil {
		return err
	}
	_, err = tree.Commit("Initial commit", &git.CommitOptions{})
	return err
}
