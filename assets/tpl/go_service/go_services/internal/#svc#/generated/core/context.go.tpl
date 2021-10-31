package core

import (
	"go.uber.org/zap"
	"{{.GoModule}}/internal/{{.ServiceName}}/app"
)

type MifyServiceContext struct {
	ServiceName    string
	Logger         *zap.Logger
	SugaredLogger  *zap.SugaredLogger
	ServiceContext app.ServiceContext
}

func NewMifyServiceContext(serviceName string) (MifyServiceContext, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return MifyServiceContext{}, err
	}

	defer logger.Sync()

	svcCtx, err := app.NewServiceContext()
	if err != nil {
		return err
	}

	context := MifyServiceContext{
		ServiceName:    serviceName,
		Logger:         logger,
		SugaredLogger:  logger.Sugar(),
		ServiceContext: svcCtx,
	}

	return context, nil
}
