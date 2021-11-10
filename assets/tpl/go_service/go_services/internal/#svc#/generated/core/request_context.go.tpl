{{- .Workspace.TplHeader}}

package core

import (
	"go.uber.org/zap"

	"context"
)

type MifyRequestContextBuilder struct {
	requestId string
	protocol  string
	urlPath   string
	logger    *zap.Logger
	metrics   *RequestMetrics

	serviceContext *MifyServiceContext
}

func NewMifyRequestContextBuilder(serviceContext *MifyServiceContext) *MifyRequestContextBuilder {
	return &MifyRequestContextBuilder{
		logger:         serviceContext.Logger(),
		serviceContext: serviceContext,
		metrics:        NewRequestMetrics(),
	}
}

func (b *MifyRequestContextBuilder) SetRequestID(requestId string) *MifyRequestContextBuilder {
	b.requestId = requestId
	return b
}

func (b *MifyRequestContextBuilder) SetProtocol(protocol string) *MifyRequestContextBuilder {
	b.protocol = protocol
	return b
}

func (b *MifyRequestContextBuilder) SetURLPath(path string) *MifyRequestContextBuilder {
	b.urlPath = path
	return b
}

func (b *MifyRequestContextBuilder) GetURLPath() string {
	return b.urlPath
}

func (b *MifyRequestContextBuilder) GetMetrics() *RequestMetrics {
	return b.metrics
}

func (b *MifyRequestContextBuilder) Logger() *zap.Logger {
	return b.logger.With(
		zap.String("request_id", b.requestId),
		zap.String("proto", b.protocol),
		zap.String("path", b.urlPath),
	)
}

func (b *MifyRequestContextBuilder) ServiceContext() *MifyServiceContext {
	return b.serviceContext
}

func (b *MifyRequestContextBuilder) Build(ctx context.Context) (*MifyRequestContext, error) {
	return &MifyRequestContext{
		requestId:      b.requestId,
		serviceContext: b.serviceContext,
		logger:         b.Logger(),
		requestCtx:     ctx,
	}, nil
}

type MifyRequestContext struct {
	requestId      string
	serviceContext *MifyServiceContext

	logger     *zap.Logger
	requestCtx context.Context
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
