package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
)

type AppCalculator struct {
	router http.Handler
	db     *sql.DB
	config Config
}

func New(config Config) *AppCalculator {
	app := &AppCalculator{
		config: config,
	}
	app.loadDB()
	app.loadRoutes()
	return app
}

func (a *AppCalculator) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.router,
	}

	err := a.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to mysql database: %w", err)
	}

	defer func() {
		if err = a.db.Close(); err != nil {
			fmt.Println("failed to close connection to mysql database:", err)
		}
	}()

	fmt.Println("Starting server")
	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to listen to server: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		fmt.Println("Shutting down server")
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
}

func (a *AppCalculator) loadDB() {
	cfg := mysql.NewConfig()
	cfg.User = "calculator_microservice"
	cfg.Passwd = "yixnhg64G0.*hafc2^"
	cfg.Net = "tcp"
	cfg.Addr = a.config.DbAddress
	cfg.DBName = "users_data"
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}
	a.db = db
}
