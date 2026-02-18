package main

import (
	"fmt"
	"context"
	"go-product-restapi/internal/app"
	"go-product-restapi/internal/config"
	"go-product-restapi/pkg/logger"
	"go-product-restapi/pkg/postgres"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	loggerFactory := logger.NewGlobalLoggerFactory()
	cfg := config.LoadConfig()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := postgres.NewClient(ctx, connString)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := app.NewApp(db, loggerFactory)

	go func() {
		if err := app.Run(cfg); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		panic(err)
	}
}
