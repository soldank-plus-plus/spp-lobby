package main

import (
	"log"

	"github.com/nedik/spp-lobby/controllers"
	"github.com/nedik/spp-lobby/initializers"
	"github.com/nedik/spp-lobby/routes"

	"github.com/gin-gonic/gin"
)

var (
	FaviconController      controllers.FaviconController
	FaviconRouteController routes.FaviconRouteController

	ServerController      controllers.ServerController
	ServerRouteController routes.ServerRouteController
)

func init() {
	FaviconController = controllers.NewFaviconController()
	FaviconRouteController = routes.NewFaviconRouteController(FaviconController)

	ServerController = controllers.NewServerController()
	ServerRouteController = routes.NewServerRouteController(ServerController)
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	server := gin.Default()
	initializers.SetupServer(server)
	router := server.Group("/")

	ServerRouteController.InitServerRoutes(router)
	FaviconRouteController.InitFaviconRoutes(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
