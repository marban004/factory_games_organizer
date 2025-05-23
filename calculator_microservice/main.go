package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
	microservicelogiccalculator "github.com/marban004/factory_games_organizer.git/microservice_logic_calculator"
)

var desiredResourceName = "reinforced_iron_plate"
var userId = 1
var altRecipies = [0]string{}
var db *sql.DB
var err error

func main() {
	cfg := mysql.NewConfig()
	cfg.User = "calculator_microservice"
	cfg.Passwd = "yixnhg64G0.*hafc2^"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "users_data"

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err.Error())
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Get("/calculate", basicHandler)

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("failed to listen to server:", err)
	}

}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	byteJSONRepresentation, err := microservicelogiccalculator.Calculate(userId, desiredResourceName, 0.5, altRecipies[:], db)
	if err != nil {
		fmt.Printf("Could not generate production tree for '%s', reason: %v \n", desiredResourceName, err)
	}
	w.Write(byteJSONRepresentation)
}
