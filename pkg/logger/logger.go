package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILogger interface {
	Debug(message string, args ...any)
	Info(message string, args ...any)
	Warn(message string, args ...any)
	Error(message string, args ...any)
	Fatal(message string, args ...any)

	Named(name string) ILogger
}

type GlobalLoggerFactory struct {
	globalLogger *Logger
	once         sync.Once
}

func NewGlobalLoggerFactory() *GlobalLoggerFactory {
	return &GlobalLoggerFactory{}
}

func (lf *GlobalLoggerFactory) GetGlobalLogger() *Logger {
	lf.once.Do(func() {
		pe := zap.NewProductionEncoderConfig()
		pe.EncodeTime = zapcore.ISO8601TimeEncoder
		consoleEncoder := zapcore.NewJSONEncoder(pe)

		core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)
		l := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

		lf.globalLogger = &Logger{
			logger: l,
		}
	})

	return lf.globalLogger
}

type Logger struct {
	logger *zap.Logger
	name   string
}

var _ ILogger = (*Logger)(nil)

func (l *Logger) Named(name string) ILogger {
	return &Logger{
		logger: l.logger.Named(name),
	}
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
	os.Exit(1)
}
