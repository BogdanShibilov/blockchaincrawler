package logger

import "go.uber.org/zap"

type SugaredLogger struct {
	*zap.SugaredLogger
}

func NewZap() *SugaredLogger {
	logger := zap.Must(zap.NewDevelopment()).Sugar()

	return &SugaredLogger{logger}
}
