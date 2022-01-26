package gencontext

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("15:04:05.000") + "]")
}

func initLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = SyslogTimeEncoder
	cfg.EncoderConfig.EncodeCaller = nil
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.Level.SetLevel(zap.ErrorLevel) // TODO: log to file

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
