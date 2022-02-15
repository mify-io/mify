package workspace

import (
	"github.com/go-git/go-git/v5"
	"github.com/mify-io/mify/pkg/workspace/mutators"
)

func InitGit(mutContext *mutators.MutatorContext) error {
	_, err := git.PlainInit(mutContext.GetDescription().BasePath, false)
	return err
}
