package app

import "go-product-restapi/pkg/postgres"

type App struct {
	DB *postgres.Client
}

func NewApp(db *postgres.Client) *App {
	return &App{
		DB: db,
	}
}

func (a *App) Close() {
	a.DB.Close()
}

func (a *App) Run() error {
	return nil
}
