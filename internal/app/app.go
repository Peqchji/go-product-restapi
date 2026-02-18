package app

import (
	"context"
	"fmt"
	"net/http"

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
	productService := productusecase.NewProductUseCase(a.globalLogger, productRepo)
	productHandler := producthandler.NewProductHandler(productService)

	echoServer := echo.New()
	echoServer.POST("/products", productHandler.CreateProduct())
	echoServer.PATCH("/products/:id", productHandler.PatchProduct())
	echoServer.GET("/products/:id", productHandler.GetProduct())
	echoServer.GET("/health", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	return echoServer.Start(fmt.Sprintf(":%s", conf.ServerPort))
}
