package logger

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILogger interface {
	Debug(message string, args ...any)
	Info(message string, args ...any)
	Warn(message string, args ...any)
	Error(message string, args ...any)
	Fatal(message string, args ...any)
}

type Logger struct {
	logger  *zap.Logger
	once    sync.Once
	filePtr *os.File
}

var _ ILogger = (*Logger)(nil)

func NewLogger(logFilePath string) (*Logger, error) {
	now := time.Now()
	logfile := path.Join(
		logFilePath,
		fmt.Sprintf("%s.log", now.Format("2006-01-02")),
	)

	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	pe := zap.NewProductionEncoderConfig()

	fileEncoder := zapcore.NewJSONEncoder(pe)
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), highPriority),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), highPriority),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))

	return &Logger{
		logger:  logger,
		filePtr: file,
	}, nil
}

func (l *Logger) Shutdown() {
	l.once.Do(func() {
		l.logger.Sync()
		l.filePtr.Close()
	})
}

func (l *Logger) Debug(message string, args ...any) {
	l.logger.Debug(message, zap.Any("args", args))
}

func (l *Logger) Info(message string, args ...any) {
	l.logger.Info(message, zap.Any("args", args))
}

func (l *Logger) Warn(message string, args ...any) {
	l.logger.Warn(message, zap.Any("args", args))
}

func (l *Logger) Error(message string, args ...any) {
	l.logger.Error(message, zap.Any("args", args))
}

func (l *Logger) Fatal(message string, args ...any) {
	l.logger.Fatal(message, zap.Any("args", args))
}
