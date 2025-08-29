package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/marban004/factory_games_organizer/application"
	_ "github.com/marban004/factory_games_organizer/docs"
)

//	@title			Dispatcher microservice
//	@version		1.0-go-to-hell
//	@description	This is a dispatcher microservice for Factory Games Organizer api. Api is responsible for creating, updating and authenicating api users, CRUD operations on database associated with the api and provides production calculator service.

//	@contact.name	Ur on your own
//	@contact.url	404
//	@contact.email	not_my@business.com

//	@license.name	You think I have a license?
//	@license.url	404

//	@host		79.175.222.18:40000
//	@BasePath	/

//	@OpenAPIDefinition(servers	= {@Server(url = "/", description = "a microservice host"), @Server(url = "/", description = "dispatcher microservice, microservices are differentiated by port number")})
//	@securityDefinitions.apikey	apiTokenAuth
//
//	@in							query
//	@name						jwt
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

// func main() {
// 	ctx := context.Background()
// 	client := &http.Client{
// 		Timeout: 10 * time.Second,
// 	}
// 	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://localhost:3000/health", nil)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	response, err := client.Do(request)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	fmt.Println(response.Body)
// }
