{{- .Workspace.TplHeader}}

package core

import (
	"go.uber.org/zap"

	"context"
)

type MifyRequestContextBuilder struct {
	requestId string
	logger    *zap.Logger
}

func NewMifyRequestContextBuilder(logger *zap.Logger) *MifyRequestContextBuilder {
	return &MifyRequestContextBuilder{logger: logger}
}

func (b *MifyRequestContextBuilder) SetRequestID(requestId string) *MifyRequestContextBuilder {
	b.requestId = requestId
	return b
}

func (b *MifyRequestContextBuilder) Logger() *zap.Logger {
	return b.logger
}

func (b *MifyRequestContextBuilder) Build(ctx context.Context, mifyServiceContext *MifyServiceContext) (*MifyRequestContext, error) {
	return &MifyRequestContext{
		requestId: b.requestId,
		serviceContext: mifyServiceContext,
		logger: b.logger,
		requestCtx: ctx,
	}, nil
}

type MifyRequestContext struct {
	requestId      string
	serviceContext *MifyServiceContext

	logger *zap.Logger
	requestCtx    context.Context
}

func (c *MifyRequestContext) Logger() *zap.Logger {
	return c.logger
}

func (c *MifyRequestContext) ServiceContext() *MifyServiceContext {
	return c.serviceContext
}

func (c *MifyRequestContext) RequestContext() context.Context {
	return c.requestCtx
}
