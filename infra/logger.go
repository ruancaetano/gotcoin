package infra

import (
	"go.uber.org/zap"
)

func GetLogger(tag string) *zap.Logger {
	config := zap.NewDevelopmentConfig()

	config.EncoderConfig.EncodeCaller = nil

	logger, _ := config.Build()

	return logger.Named(tag)
}
