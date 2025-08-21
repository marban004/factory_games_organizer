package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/marban004/factory_games_organizer/application"
	_ "github.com/marban004/factory_games_organizer/docs"
)

//	@title			Calculator microservice
//	@version		1.0-go-to-hell
//	@description	This is a calculator microservice for Factory Games Organizer api. It is not meant to be accessed directly. Access to the microservice should be done through dispatcher microservice.

//	@contact.name	Ur on your own
//	@contact.url	404
//	@contact.email	not_my@business.com

//	@license.name	You think I have a license?
//	@license.url	404

//	@host		192.168.100.16:8080
//	@BasePath	/

//	@OpenAPIDefinition(servers	= {@Server(url = "/", description = "a microservice host"), @Server(url = "/", description = "calculator microservice, microservices are differentiated by port number")})
//
// host is WAN address of router, need to set up port forwarding to redirect to LAN address of my laptop, also set the address to be static on my laptop, so port forwarding always goes to it.
func main() {
	app := application.New(application.LoadConfig())
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
	cancel()
}
