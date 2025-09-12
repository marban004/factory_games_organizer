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

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/marban004/factory_games_organizer/application"
	_ "github.com/marban004/factory_games_organizer/docs"
)

//	@title			Users microservice
//	@version		1.0
//	@description	This is a Users microservice for Factory Games Organizer api. It is not meant to be accessed directly. Access to the microservice should be done through dispatcher microservice.
//
//	@contact.name	Marek Banaś
//	@contact.email	marek.banas004@gmail.com
//
//	@license.name	GPL-3.0
//	@license.url	https://www.gnu.org/licenses/gpl-3.0.html
//
//	@host		79.175.222.18:8082
//	@BasePath	/
//
//	@OpenAPIDefinition(servers	= {@Server(url = "/", description = "a microservice host"), @Server(url = "/", description = "users microservice, microservices are differentiated by port number")})
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
