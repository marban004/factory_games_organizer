package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

type AppCalculator struct {
	router http.Handler
	cfg    *mysql.Config
}

func New() *AppCalculator {
	app := &AppCalculator{
		router: loadRoutes(),
		cfg:    loadConfig(),
	}
	return app
}

func (a *AppCalculator) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to listen to server: %w", err)
	}
	return nil
}

func loadConfig() *mysql.Config {
	cfg := mysql.NewConfig()
	cfg.User = "calculator_microservice"
	cfg.Passwd = "yixnhg64G0.*hafc2^"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "users_data"
	return cfg
}
