package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

type AppUsers struct {
	router     http.Handler
	db         *sql.DB
	secret     []byte
	config     Config
	certPath   string
	secretPath string
}

func New(config Config) *AppUsers {
	app := &AppUsers{
		config: config,
	}
	app.secretPath = config.ServerSecretPath
	app.certPath = config.ServerCertPath
	app.loadSecret()
	app.loadDB()
	app.loadRoutes()
	return app
}

func (a *AppUsers) Start(ctx context.Context) error {
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
		err = server.ListenAndServeTLS(a.certPath, a.secretPath)
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

func (a *AppUsers) loadDB() {
	cfg := mysql.NewConfig()
	cfg.User = "users_microservice"
	cfg.Passwd = "bxu7%^yhag##KKL"
	cfg.Net = "tcp"
	cfg.Addr = a.config.DbAddress
	cfg.DBName = "users"
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}
	a.db = db
}

// implement reading file designated by the path in a.config.ServerCertPath
func (a *AppUsers) loadSecret() {
	fileContents, err := os.ReadFile(a.config.ServerSecretPath)
	if err != nil {
		panic("could not open server secret key file")
	}
	a.secret = fileContents
}
