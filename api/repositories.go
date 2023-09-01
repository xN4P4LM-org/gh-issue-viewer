package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// Check if repository is valid and/or accessible
//

// function to check if the repository is valid and/or accessible
func checkRepository(request_context *gin.Context) {

	type Status struct {
		Valid   bool   `json:"valid"`
		Message string `json:"message"`
		Error   string `json:"error"`
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

	owner_name := request_context.Query("owner")

	if owner_name == "" {

		response := Status{
			Valid:   false,
			Message: "owner is required",
			Error:   err.Error(),
		}

		request_context.JSON(
			http.StatusBadRequest,
			gin.H{"error": response})
	}

	repository_name := request_context.Query("repository")

	if repository_name == "" {

		response := Status{
			Valid:   false,
			Message: "repository is required",
			Error:   err.Error(),
		}

		request_context.JSON(
			http.StatusBadRequest,
			gin.H{"error": response})
	}

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

		request_context.JSON(http.StatusOK, gin.H{"Repository": response, "Request Headers": headers})
	}

}
