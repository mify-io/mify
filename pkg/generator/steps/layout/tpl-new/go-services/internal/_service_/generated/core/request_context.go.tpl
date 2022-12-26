{{- .TplHeader}}
// vim: set ft=go:

package core

import (
	"go.uber.org/zap"

	"context"
	"net/http"
	"time"

	"{{.MetricsImportPath}}"
)

type MifyRequestContextBuilder struct {
	requestId string
	protocol  string
	urlPath   string
	logger    *zap.Logger
	metrics   *metrics.RequestMetrics

	serviceContext *MifyServiceContext
}

func NewMifyRequestContextBuilder(serviceContext *MifyServiceContext) *MifyRequestContextBuilder {
	return &MifyRequestContextBuilder{
		logger:         serviceContext.Logger(),
		serviceContext: serviceContext,
		metrics:        metrics.NewRequestMetrics(),
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

func (b *MifyRequestContextBuilder) GetMetrics() *metrics.RequestMetrics {
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

func (b *MifyRequestContextBuilder) Build(r *http.Request, rw http.ResponseWriter) (*MifyRequestContext, error) {
	return &MifyRequestContext{
		MifyServiceContext: b.serviceContext,

		requestId:      b.requestId,
		logger:         b.Logger(),
		request:        r,
		responseWriter: rw,
	}, nil
}

type MifyRequestContext struct {
	*MifyServiceContext

	requestId      string
	logger         *zap.Logger
	request        *http.Request
	responseWriter http.ResponseWriter
}

// WithGoContext returns a shallow copy of MifyRequestContext which allows
// to pass different go context into nested functions, for instance to add
// custom deadlines.
func (c *MifyRequestContext) WithGoContext(goCtx context.Context) *MifyRequestContext {
	return &MifyRequestContext{
		MifyServiceContext: c.MifyServiceContext,
		requestId:      c.requestId,
		logger:         c.logger,
		request:        c.request.WithContext(goCtx),
		responseWriter: c.responseWriter,
	}
}

func (c *MifyRequestContext) Logger() *zap.Logger {
	return c.logger
}

func (c *MifyRequestContext) Request() *http.Request {
	return c.request
}

func (c *MifyRequestContext) RequestContext() context.Context {
	return c.request.Context()
}

func (c *MifyRequestContext) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

// context.Context

func (c *MifyRequestContext) Deadline() (deadline time.Time, ok bool) {
	return c.request.Context().Deadline()
}

func (c *MifyRequestContext) Done() <-chan struct{} {
	return c.request.Context().Done()
}

func (c *MifyRequestContext) Err() error {
	return c.request.Context().Err()
}

func (c *MifyRequestContext) Value(key interface{}) interface{} {
	return c.request.Context().Value(key)
}
