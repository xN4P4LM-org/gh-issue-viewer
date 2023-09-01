package api

import (
	"context"

	"github.com/gin-gonic/gin"
)

// global context
var CTX context.Context

// setup the web api
func SetupApi(web *gin.Engine, ctx context.Context) {

	// set the global context
	CTX = ctx

	// define the api prefix
	api_prefix := web.Group("/api")

	// route to check if the repository is valid and/or accessible
	api_prefix.GET("repository/valid", checkRepository)

	// route to get all issues for a given repository
	api_prefix.GET("repository/issues", getIssues)

}
