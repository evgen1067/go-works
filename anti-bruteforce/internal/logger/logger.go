package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/evgen1067/anti-bruteforce/internal/common"
	"github.com/evgen1067/anti-bruteforce/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultOutputPaths = []string{"out.log"}

type Logger struct {
	logger *zap.Logger
}

func NewLogger(cfg *config.Config) (*Logger, error) {
	var level zap.AtomicLevel
	switch cfg.Logger.Level {
	case "error":
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "debug":
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		return nil, common.ErrUndefinedLoggerLevel
	}

	var file []string
	if cfg.Logger.File == "" {
		file = defaultOutputPaths
	} else {
		file = append(file, cfg.Logger.File)
	}
	pathDir := "logs"
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		err := os.Mkdir(pathDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	for i := range file {
		file[i] = filepath.Join(pathDir, file[i])
	}
	// file = append(file, "stdout")

	zapCfg := zap.Config{
		Level:       level,
		Encoding:    "console",
		OutputPaths: file,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			TimeKey:     "time",
			EncodeLevel: CustomEncodeLevel,
			EncodeTime:  CustomEncodeTime,
		},
	}

	zapLogger, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{
		logger: zapLogger,
	}, nil
}

func CustomEncodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + l.String() + "]")
}

func CustomEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2/Jan/2006:15:04:05 -0700"))
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}
