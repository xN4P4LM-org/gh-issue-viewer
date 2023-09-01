package frontend

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// global context
var CTX context.Context

// setup the front end
func SetupFrontend(web *gin.Engine, ctx context.Context) {

	// set the global context
	CTX = ctx

	// the directory containing the scripts
	scriptsDir := "frontend/static/scripts/"

	var individualScripts []string

	// if production mode is enabled, combine the scripts into one file
	// otherwise, serve the scripts individually
	if !gin.IsDebugging() {
		// combine the scripts
		combineJS(scriptsDir)
	} else {
		// serve the scripts individually
		individualScripts = getAllScripts()
	}

	// serve the favicon statically
	web.Static("img", "./frontend/static/img")

	// serve the css statically
	web.Static("css", "./frontend/static/css")

	// serve the scripts statically
	web.Static("scripts", scriptsDir)

	// load the HTML templates from frontend/templates/**/*
	web.LoadHTMLGlob("frontend/templates/**/*")

	if !gin.IsDebugging() {
		// setup the router
		Router(web, []string{""})
	} else {
		// setup the router with the scripts
		Router(web, individualScripts)
	}

}

func Router(web *gin.Engine, scripts []string) {
	web.GET("/", func(c *gin.Context) {
		index(c, scripts)
	})
	web.GET("/index", func(c *gin.Context) {
		index(c, scripts)
	})
	web.GET("/index.html", func(c *gin.Context) {
		index(c, scripts)
	})
}

func index(request_context *gin.Context, scripts []string) {

	request_context.HTML(
		http.StatusOK,
		"views/index.html",
		gin.H{
			"title":   "gh-issue-viewer",
			"scripts": scripts,
			"debug":   gin.IsDebugging(),
		})
}
