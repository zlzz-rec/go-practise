package log

import "go.uber.org/zap"

var Logger *zap.Logger

func Init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
}

