package app

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v5"

	"go-product-restapi/internal/config"
	"go-product-restapi/internal/handler/producthandler"
	"go-product-restapi/internal/repositories/productrepository"
	"go-product-restapi/internal/usecases/productusecase"
	"go-product-restapi/pkg/logger"
	"go-product-restapi/pkg/postgres"
)

type App struct {
	globalLogger logger.ILogger
	DB           *postgres.Client
}

func NewApp(
	db *postgres.Client,
	loggerFactory *logger.GlobalLoggerFactory,
) *App {
	return &App{
		DB:           db,
		globalLogger: loggerFactory.GetGlobalLogger(),
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	done := make(chan struct{})
	var err error

	go func() {
		defer close(done)

		if a.DB != nil {
			a.DB.Close()
		}

		if a.globalLogger != nil {
			err = a.globalLogger.Shutdown()
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return err
	}
}

func (a *App) Run(conf *config.Config) error {
	db := a.DB.Pool

	productMapper := productrepository.NewPostgresProductMapper()
	productRepo := productrepository.NewPostgresProductRepository(db, a.globalLogger, productMapper)
	productService := productusecase.NewProductUseCase(productRepo)
	productHandler := producthandler.NewProductHandler(productService)

	echo := echo.New()
	echo.POST("/products", productHandler.CreateProduct())
	echo.PATCH("/products/:id", productHandler.PatchProduct())

	return echo.Start(fmt.Sprintf(":%s", conf.ServerPort))
}
