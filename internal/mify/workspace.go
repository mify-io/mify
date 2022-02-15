package mify

import (
	"fmt"
	"path"

	"github.com/mify-io/mify/internal/mify/util"
	"github.com/mify-io/mify/pkg/workspace/mutators"
	"github.com/mify-io/mify/pkg/workspace/mutators/workspace"
)

var vcsTemplates = []string{"none", "git"}

func CreateWorkspace(ctx *CliContext, parentDir string, name string, vcs string) error {
	mutCtx := mutators.NewMutatorContext(ctx.Ctx, ctx.Logger, nil)

	err := workspace.CreateWorkspace(mutCtx, parentDir, name)
	if err != nil {
		return err
	}

	ctx.WorkspacePath = path.Join(parentDir, name) // TODO: remove this hack
	err = ctx.InitWorkspaceDescription()
	if err != nil {
		return err
	}
	mutCtx = mutators.NewMutatorContext(ctx.Ctx, ctx.Logger, ctx.MustGetWorkspaceDescription())

	err = util.ValidateStrArg(vcs, vcsTemplates)
	if err != nil {
		return fmt.Errorf("invalid vcs template: %w", err)
	}

	switch vcs {
	case "git":
		err = workspace.InitGit(mutCtx)
		if err != nil {
			return err
		}
	}

	return nil
}
