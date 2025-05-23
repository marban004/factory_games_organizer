package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"

	"github.com/marban004/factory_games_organizer/calculator_microservice/application"
)

var desiredResourceName = "reinforced_iron_plate"
var userId = 1
var altRecipies = [0]string{}
var db *sql.DB

func main() {
	app := application.New()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
	cancel()
}
