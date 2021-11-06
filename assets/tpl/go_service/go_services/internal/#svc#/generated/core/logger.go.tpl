package core

import (
	"sync"

	"go.uber.org/zap"
)

type MifyLoggerWrapper struct {
	logger     *zap.Logger
	loggersFor map[string]*zap.Logger
	rwMutex    sync.RWMutex
}

func NewMifyLoggerWrapper(ctx *MifyServiceContext) (*MifyLoggerWrapper, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return &MifyLoggerWrapper{}, err
	}

	logger = logger.With(
		zap.String("service_name", ctx.serviceName),
		zap.String("hostname", ctx.hostname),
	)

	defer logger.Sync()

	return &MifyLoggerWrapper{
		logger:     logger,
		loggersFor: make(map[string]*zap.Logger),
	}, nil
}

func (w *MifyLoggerWrapper) Logger() *zap.Logger {
	return w.logger
}

func (w *MifyLoggerWrapper) LoggerFor(component string) *zap.Logger {
	w.rwMutex.RLock()
	loggerFor, ok := w.loggersFor[component]
	w.rwMutex.RUnlock()

	if ok {
		return loggerFor
	}

	w.rwMutex.Lock()
	defer w.rwMutex.Unlock()

	loggerFor, ok = w.loggersFor[component]
	if ok {
		return loggerFor
	}

	loggerFor = w.logger.With(zap.String("component", component))
	w.loggersFor[component] = loggerFor

	return loggerFor
}
