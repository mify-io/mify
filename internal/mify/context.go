package mify

import (
	"context"
	"log"
	"os"

	"github.com/chebykinn/mify/pkg/workspace"
	"github.com/chebykinn/mify/pkg/workspace/mutators"
)

type CliContext struct {
	Logger *log.Logger
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewContext() *CliContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &CliContext{
		Logger: log.New(os.Stdout, "", 0),
		Ctx:    ctx,
		Cancel: cancel,
	}
}

func initMutatorCtx(ctx *CliContext, basePath string) (*mutators.MutatorContext, error) {
	descr, err := workspace.InitDescription(basePath)
	if err != nil {
		return nil, err
	}

	return mutators.NewMutatorContext(ctx.Ctx, ctx.Logger, &descr), nil
}
