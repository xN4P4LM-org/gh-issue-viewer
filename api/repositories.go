package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// Check if repository is valid and/or accessible
//

// function to check if the repository is valid and/or accessible
func checkRepository(request_context *gin.Context) {

	type RequestData struct {
		Owner      string `json:"owner"`
		Repository string `json:"repository"`
	}

	github_api, err := setupGithubApi()

	if err != nil {

		// create the json response
		response := Status{
			Valid:   false,
			Message: "Error setting up github api",
			Error:   err.Error(),
		}

		request_context.JSON(
			http.StatusInternalServerError,

			gin.H{"error": response})

		return
	}

	// Read the raw request body
	body, err := io.ReadAll(request_context.Request.Body)
	if err != nil {
		// Handle the error
		return
	}

	// Create an instance of the struct to hold the decoded JSON data
	requestData := RequestData{}

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		// Handle the error
		return
	}

	// Access the owner field from the decoded JSON data
	owner_name := requestData.Owner
	repository_name := requestData.Repository

	_, response, request_err := github_api.Repositories.Get(CTX, owner_name, repository_name)

	// capture the headers
	headers := Headers{
		RateLimitRemaining: response.Rate.Remaining,
		RequestID:          response.Header.Get("X-GitHub-Request-Id"),
		Date:               response.Header.Get("Date"),
	}

	// confirm there was no error
	if request_err != nil {

		response := Status{
			Valid:   false,
			Message: "Repository not Found",
			Error:   request_err.Error(),
		}

		// return the error in a json response
		request_context.JSON(http.StatusInternalServerError, gin.H{"error": response})
	} else {

		response := Status{
			Valid:   true,
			Message: "repository is valid and accessible",
			Error:   "",
		}

		request_context.JSON(http.StatusOK, gin.H{"repository": response, "Request Headers": headers})
	}

}
