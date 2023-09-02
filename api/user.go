package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func checkUser(request_context *gin.Context) {

	type RequestData struct {
		User string `json:"user"`
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
	user_name := requestData.User

	if user_name == "" {

		response := Status{
			Valid:   false,
			Message: "user is required",
			Error:   err.Error(),
		}

		request_context.JSON(
			http.StatusBadRequest,
			gin.H{"error": response})
	}

	_, _, err_user := github_api.Users.Get(CTX, user_name)

	if err_user != nil {

		response := Status{
			Valid:   false,
			Message: "user is invalid",
			Error:   err.Error(),
		}

		request_context.JSON(
			http.StatusBadRequest,
			gin.H{"error": response})
	}

	response := Status{
		Valid:   true,
		Message: "user is valid",
		Error:   "",
	}

	request_context.JSON(
		http.StatusOK,
		gin.H{"user": response})

}
