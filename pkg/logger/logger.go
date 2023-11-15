package logger

import "go.uber.org/zap"

type ZapLogger struct {
	*zap.SugaredLogger
}

func NewZap() *ZapLogger {
	logger := zap.Must(zap.NewDevelopment()).Sugar()

	return &ZapLogger{logger}
}
