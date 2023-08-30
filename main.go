package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/xn4p4lm/gh-issue-viewer/api"
	"github.com/xn4p4lm/gh-issue-viewer/frontend"

	"github.com/gin-gonic/gin"
)

// ping function that returns a JSON response with the
// provided status code (200) and the response pong
func ping(request_context *gin.Context) {
	request_context.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"status":  http.StatusOK})
}

func setupRouter() *gin.Engine {

	// setup the web server
	web := gin.Default()

	// Ping test
	web.GET("api/ping", ping)

	// log that the server has been initialized
	log.Print("Server initialized")

	return web
}

func main() {
	// setup the context
	ctx := context.Background()

	// setup the web server
	web := setupRouter()

	// setup the API
	api.SetupApi(web, ctx)
	log.Print("API initialized")

	// setup the frontend
	frontend.SetupFrontend(web, ctx)
	log.Print("Frontend initialized")

	// check if gin is running in production
	if gin.Mode() == "release" {

		// get the port from environment variable
		port := ":" + os.Getenv("PORT")

		// run the server
		web.Run(port)
		// log that the server has been started on the domain and port
		log.Print("Starting Server @ localhost:" + port)

	}

	if gin.Mode() == "debug" {

		// Listen and Server in 0.0.0.0:8080
		web.Run(":8080")

	}
}
