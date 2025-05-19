package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func createLogger() (logger *zap.Logger) {
	config := zap.NewDevelopmentConfig()
	// 设置日志格式为控制台格式
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.Encoding = "console"

	logger, _ = config.Build()
	return logger
}

func init() {
	Logger = createLogger()
}
