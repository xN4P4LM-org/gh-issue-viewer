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

	// if environment variable ENVIRONMENT is production, then run the server using the FQDN and PORT
	if os.Getenv("ENVIRONMENT") == "production" {

		// get the FQDN from environment variable
		domain := os.Getenv("FQDN")

		// get the port from environment variable
		port := os.Getenv("PORT")

		// create the domain and port string
		domain_port := domain + ":" + port

		// run the server
		web.Run(domain_port)

	}

	if os.Getenv("ENVIRONMENT") != "production" {

		// Listen and Server in 0.0.0.0:8080
		web.Run(":8080")

	}
}
