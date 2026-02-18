package logger_test

import (
	"testing"

	"go-product-restapi/pkg/logger"
)

func TestGlobalLogger(t *testing.T) {
	factory := logger.NewGlobalLoggerFactory()
	log := factory.GetGlobalLogger()

	if log == nil {
		t.Fatal("Expected global logger to be non-nil")
	}

	log.Debug("debug message")
	log.Info("info message")
	log.Warn("warn message")
	log.Error("error message")

	named := log.Named("test-component")
	named.Info("named logger message")
}