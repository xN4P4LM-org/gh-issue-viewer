package frontend

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

// global context
var CTX context.Context

// setup the front end
func SetupFrontend(web *gin.Engine, ctx context.Context) {

	// set the global context
	CTX = ctx

	// combine the scripts
	scripts_path := combineJS()

	// serve the favicon statically
	web.Static("favicon.ico", "./frontend/static/img/favicon.ico")

	// serve the css statically
	web.Static("css", "./frontend/static/css")

	// serve the scripts statically
	web.Static("scripts", scripts_path)

	// setup the functions for the templates
	web.SetFuncMap(template.FuncMap{
		"formatDate": formatDate,
	})

	// load the HTML templates from frontend/templates/**/*
	web.LoadHTMLGlob("frontend/templates/**/*")

	// setup the router
	Router(web)

}

// setup the router
func Router(web *gin.Engine) {
	web.GET("/", index)
	web.GET("/index", index)
	web.GET("/index.html", index)
}

// function to render the time page
func formatDate() time.Time {
	return time.Now()
}

// function to render the index page
func index(request_context *gin.Context) {

	// render the index page
	request_context.HTML(
		http.StatusOK,
		"views/index.html",
		gin.H{
			"title": "gh-issue-viewer",
			"now":   time.Now(),
		})
}

// function to combine the javascript files
func combineJS() string {
	// the directory containing the scripts
	scriptsDir := "frontend/static/scripts/"

	// combined file name
	outputFile := "combined.js"

	// the output file
	outputFilePath := scriptsDir + outputFile

	// read the directory
	files, err := os.ReadDir(scriptsDir)

	// if there was an error reading the directory, exit the program
	if err != nil {
		fmt.Println("Error reading directory:", err)
		panic(err)
	}

	// create a string to hold the combined content
	var combinedContent string

	// loop through the files
	for _, file := range files {
		// if file is scripts_combined.js, skip it
		if file.Name() == outputFile {
			continue
		}

		// if file is a .js file, combine it
		if filepath.Ext(file.Name()) == ".js" {

			// read the file
			content, err := os.ReadFile(filepath.Join(scriptsDir, file.Name()))

			// if there was an error reading the file, exit the program
			if err != nil {
				fmt.Println("Error reading file:", err)
				panic(err)
			}

			// add the content to the combined content
			combinedContent += string(content) + "\n"
		}
	}

	// setup the tool to minify the content
	minify := minify.New()

	minify.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

	// minify the content
	combinedContent, err = minify.String("text/javascript", combinedContent)
	if err != nil {
		panic(err)
	}

	// write the combined content to the output file
	err = os.WriteFile(outputFilePath, []byte(combinedContent), 0644)

	// if there was an error writing to the output file, exit the program
	if err != nil {
		fmt.Println("Error writing to combined file:", err)
		panic(err)
	}

	// return the output file
	return scriptsDir
}
