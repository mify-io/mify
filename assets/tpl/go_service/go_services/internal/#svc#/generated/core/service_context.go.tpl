{{- .Workspace.TplHeader}}

package core

import (
	"os"

	"go.uber.org/zap"
	"{{.GoModule}}/internal/{{.ServiceName}}/app"
)

type MifyServiceContext struct {
	serviceName string
	hostname    string

	logger        *zap.Logger

	serviceContext *app.ServiceContext
}

func NewMifyServiceContext(serviceName string) (*MifyServiceContext, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return &MifyServiceContext{}, err
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return &MifyServiceContext{}, err
	}

	logger = logger.With(
		zap.String("service_name", serviceName),
		zap.String("hostname", hostname),
	)

	defer logger.Sync()

	svcCtx, err := app.NewServiceContext()
	if err != nil {
		return &MifyServiceContext{}, err
	}

	context := &MifyServiceContext{
		serviceName:    serviceName,
		hostname:       hostname,
		logger:         logger,
		serviceContext: svcCtx,
	}

	return context, nil
}

func (c *MifyServiceContext) ServiceName() string {
	return c.serviceName
}

func (c *MifyServiceContext) Hostname() string {
	return c.hostname
}

func (c *MifyServiceContext) Logger() *zap.Logger {
	return c.logger
}

func (c *MifyServiceContext) ServiceContext() *app.ServiceContext {
	return c.serviceContext
}
