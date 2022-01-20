package mutators

import (
	"context"
	"log"

	"github.com/chebykinn/mify/pkg/workspace"
)

type MutatorContext struct {
	logger      *log.Logger
	ctx         context.Context
	cancel      context.CancelFunc
	description *workspace.Description // Description can be null when workspace is not created yet
}

func NewMutatorContext(ctx context.Context, log *log.Logger, description *workspace.Description) *MutatorContext {
	ctx, cancel := context.WithCancel(ctx)
	return &MutatorContext{
		logger:      log,
		ctx:         ctx,
		cancel:      cancel,
		description: description,
	}
}

func (c *MutatorContext) GetLogger() *log.Logger {
	if c.logger == nil {
		panic("Logger is not filled")
	}

	return c.logger
}

func (c *MutatorContext) GetCtx() context.Context {
	if c.ctx == nil {
		panic("Ctx is not filled")
	}

	return c.ctx
}

func (c *MutatorContext) GetCancel() context.CancelFunc {
	if c.cancel == nil {
		panic("Cancel is not filled")
	}

	return c.cancel
}

func (c *MutatorContext) GetDescription() *workspace.Description {
	if c.description == nil {
		panic("Description is not filled")
	}

	return c.description
}
