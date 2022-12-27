{{- .Workspace.TplHeader}}
// vim: set ft=go:

package logs

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MifyLoggerWrapper struct {
	logger     *zap.Logger
	loggersFor map[string]*zap.Logger
	rwMutex    sync.RWMutex
}

type MifyServiceContext interface {
	ServiceName() string
	Hostname()    string
}

func NewMifyLoggerWrapper(ctx MifyServiceContext) (*MifyLoggerWrapper, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "@timestamp"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	logger, err := config.Build(zap.Fields(
		zap.String("service_name", ctx.ServiceName()),
		zap.String("hostname", ctx.Hostname()),
	))
	if err != nil {
		return &MifyLoggerWrapper{}, err
	}

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
