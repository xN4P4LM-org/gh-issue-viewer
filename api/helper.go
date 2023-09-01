package api

import (
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v54/github"
)

//
// Github
//

type Headers struct {
	RateLimitRemaining int    `json:"X-RateLimit-Remaining"`
	RequestID          string `json:"X-GitHub-Request-Id"`
	Date               string `json:"Date"`
}

// function to setup the github api
func setupGithubApi() (*github.Client, error) {

	// Shared transport to reuse TCP connections.
	tr := http.DefaultTransport

	// get the github app id from the environment as int64
	app_id, err_app_id := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)

	if err_app_id != nil {
		log.Panic(err_app_id)
	}

	if app_id == 0 {
		log.Panic("GITHUB_APP_ID is empty")
	}

	// get the installation id from the environment as int64
	installation_id, err_installation_id := strconv.ParseInt(os.Getenv("GITHUB_APP_INSTALLATION_ID"), 10, 64)

	if err_installation_id != nil {
		log.Panic(err_installation_id)
	}

	if installation_id == 0 {
		log.Panic("GITHUB_APP_INSTALLATION_ID is empty")
	}

	var private_key []byte

	// if gin is running in production mode, use the github app private key from the environment and save it to a file
	if gin.Mode() == "release" {
		// get the private key from the environment
		private_key_env := os.Getenv("GITHUB_APP_PRIVATE_KEY")

		// if the private key is empty, return an error and panic
		if private_key_env == "" {
			log.Panic("GITHUB_APP_PRIVATE_KEY is empty")
		}

		// save the private key to a byte array
		private_key = []byte(private_key_env)

	} else {

		var read_key_err error

		// load the private key from the file
		private_key, read_key_err = os.ReadFile(path.Join("secure", "private-key.pem"))

		// check for errors
		if read_key_err != nil {
			log.Panic(read_key_err)
		}
	}

	// Wrap the shared transport for use with the app ID 1 authenticating with installation ID 99.
	itr, err := ghinstallation.New(
		tr,              // http.transport
		app_id,          // app_id
		installation_id, // installation_id
		private_key)     // private_key_file
	if err != nil {
		log.Fatal(err)
	}

	// Use installation transport with github.com/google/go-github
	github_api := github.NewClient(&http.Client{Transport: itr})

	return github_api, nil

}
