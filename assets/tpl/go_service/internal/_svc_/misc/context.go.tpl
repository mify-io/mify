package misc

import (
	"go.uber.org/zap"
)

type MifyContext struct {
	ServiceName   string
	Logger        *zap.Logger
	SugaredLogger *zap.SugaredLogger
}

func initContext(serviceName string) (MifyContext, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return MifyContext{}, err
	}

	defer logger.Sync()

	context := MifyContext{
		ServiceName:   serviceName,
		Logger:        logger,
		SugaredLogger: logger.Sugar(),
	}

	return context, nil
}
