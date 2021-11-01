{{- .Workspace.TplHeader}}

package core

import (
	"os"

	"go.uber.org/zap"
	"{{.GoModule}}/internal/{{.ServiceName}}/app"
)

type MifyServiceContext struct {
	ServiceName string
	Hostname    string

	Logger        *zap.Logger
	SugaredLogger *zap.SugaredLogger

	ServiceContext app.ServiceContext
}

func NewMifyServiceContext(serviceName string) (MifyServiceContext, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return MifyServiceContext{}, err
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return MifyServiceContext{}, err
	}

	logger = logger.With(
		zap.String("service_name", serviceName),
		zap.String("hostname", hostname),
	)

	defer logger.Sync()

	svcCtx, err := app.NewServiceContext()
	if err != nil {
		return MifyServiceContext{}, err
	}

	context := MifyServiceContext{
		ServiceName:    serviceName,
		Hostname:       hostname,
		Logger:         logger,
		SugaredLogger:  logger.Sugar(),
		ServiceContext: svcCtx,
	}

	return context, nil
}
