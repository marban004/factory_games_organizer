//     This is Factory Games Organizer api. Api is responsible for creating, updating and authenicating api users, CRUD operations on database associated with the api and provides production calculator service.
//     Copyright (C) 2025  Marek Banaś

//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.

//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.

//     You should have received a copy of the GNU General Public License
//     along with this program.  If not, see https://www.gnu.org/licenses/.

package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	custommiddleware "github.com/marban004/factory_games_organizer/custom_middleware"
)

type AppCalculator struct {
	router      http.Handler
	db          *sql.DB
	config      Config
	statTracker *custommiddleware.DefaultApiStatTracker
}

func New(config Config) *AppCalculator {
	app := &AppCalculator{
		config: config,
	}
	app.statTracker = &custommiddleware.DefaultApiStatTracker{MaxLen: config.TrackerCapacity, Period: config.TrackerTimePeriod, ApiStatsFile: config.ApiStatsFile, DumpStats: config.DumpStats}
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
		err = server.ListenAndServeTLS(a.config.ServerCertPath, a.config.ServerSecretPath)
		if err != nil {
			ch <- fmt.Errorf("failed to listen to server: %w", err)
		}
		close(ch)
	}()
	a.statTracker.StartTracker(ctx)

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
