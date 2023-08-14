package service

import (
	"github.com/evgen1067/anti-bruteforce/internal/logger"
	"go.uber.org/zap"
)

type LoggerService struct {
	logger *logger.Logger
}

func NewLogger(logger *logger.Logger) *LoggerService {
	return &LoggerService{
		logger: logger,
	}
}

func (l *LoggerService) Error(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *LoggerService) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *LoggerService) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *LoggerService) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}
