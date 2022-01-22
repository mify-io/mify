package mify

import (
	"github.com/mify-io/mify/pkg/workspace/mutators"
	"github.com/mify-io/mify/pkg/workspace/mutators/workspace"
)

func CreateWorkspace(ctx *CliContext, basePath string, name string) error {
	mutCtx := mutators.NewMutatorContext(ctx.Ctx, ctx.Logger, nil)

	err := workspace.CreateWorkspace(mutCtx, basePath, name)
	if err != nil {
		return err
	}

	return nil
}
