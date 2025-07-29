package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/marban004/factory_games_organizer/application"
)

func main() {
	app := application.New(application.LoadConfig())
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
	cancel()
}

// ctx context.Context, db *sql.DB, input prototypes.JSONInput, update prototypes.JSONInput, ids []int
