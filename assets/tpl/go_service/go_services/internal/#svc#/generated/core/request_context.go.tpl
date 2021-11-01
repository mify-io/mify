{{- .Workspace.TplHeader}}

package core

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MifyRequestContext struct {
	RequestId      uuid.UUID
	ServiceContext MifyServiceContext

	Logger        *zap.Logger
	SugaredLogger *zap.SugaredLogger
}

func NewMifyRequestContext(mifyServiceContext MifyServiceContext) (MifyRequestContext, error) {
	requestId := uuid.New()
	logger := mifyServiceContext.Logger.With(
		zap.String("request_id", requestId.String()),
	)

	context := MifyRequestContext{
		RequestId:      requestId,
		ServiceContext: mifyServiceContext,
		Logger:         logger,
		SugaredLogger:  logger.Sugar(),
	}
	return context, nil
}
