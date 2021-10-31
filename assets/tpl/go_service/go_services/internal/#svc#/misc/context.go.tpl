package misc

import (
	"go.uber.org/zap"
)

type MifyServiceContext struct {
	ServiceName   string
	Logger        *zap.Logger
	SugaredLogger *zap.SugaredLogger
}

func NewMifyServiceContext(serviceName string) (MifyServiceContext, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return MifyServiceContext{}, err
	}

	defer logger.Sync()

	context := MifyServiceContext{
		ServiceName:   serviceName,
		Logger:        logger,
		SugaredLogger: logger.Sugar(),
	}

	return context, nil
}
