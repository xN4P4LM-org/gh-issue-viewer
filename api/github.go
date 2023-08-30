package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v54/github"
	"golang.org/x/oauth2"
)

// global context
var CTX context.Context

// setup the web api
func SetupApi(web *gin.Engine, ctx context.Context) {

	// set the global context
	CTX = ctx

	// define the api prefix
	api_prefix := web.Group("/api")

	// route to get the username for the passed api token
	api_prefix.GET("user", getUser)

	// route to get all issues for a given repository
	api_prefix.GET("issues", getIssues)

}

//
// Github
//

// function to setup the github api
func setupGithubApi(token string) *github.Client {

	// create the bearer token by stripping the "Bearer " prefix
	bearer_token := strings.TrimPrefix(token, "Bearer ")

	// configure the oauth2 token source
	ts := oauth2.StaticTokenSource(

		&oauth2.Token{AccessToken: bearer_token},
	)

	// configure the oauth2 token
	tc := oauth2.NewClient(CTX, ts)

	// configure the github client
	github_api := github.NewClient(tc)

	// set the user agent to the name of the application
	github_api.UserAgent = "gh-issue-viewer"

	// return the github api
	return github_api
}

//
// User
//

// function to get the username for the passed api token
func getUser(request_context *gin.Context) {

	// get the token from the authentication header
	token := request_context.GetHeader("Authorization")

	// get the username for the token
	user, status_code, err := getUsername(token)

	// confirm there was no error
	if err != nil {

		// return the error in a json response
		request_context.JSON(status_code, gin.H{"error": err.Error()})
		return
	}

	// get the username
	username := user.GetLogin()

	// get the user id
	userID := user.GetID()

	// get the user avatar url
	userAvatarUrl := user.GetAvatarURL()

	// return the username in a json response
	request_context.JSON(status_code, gin.H{
		"username":        username,
		"user_id":         userID,
		"user_avatar_url": userAvatarUrl})

}

// function to get the username for the passed api token
func getUsername(token string) (*github.User, int, error) {

	// setup the github api
	github_api := setupGithubApi(token)

	// get the current user
	user, response, err := github_api.Users.Get(CTX, "")

	// confirm there was no error
	if err != nil {
		return nil, response.StatusCode, err
	}

	// get the username

	// return the username
	return user, response.StatusCode, nil

}

//
// Issues
//

// function to get all issues for a given repository
func getIssues(request_context *gin.Context) {

	// define the assignee structure
	type Assignee struct {
		Login      string `json:"login"`
		Html_Url   string `json:"html_url"`
		Avatar_Url string `json:"avatar_url"`
	}

	// define the issue structure
	type Issue struct {
		Title          string     `json:"title"`
		Body           string     `json:"body"`
		Number         int        `json:"number"`
		State          string     `json:"state"`
		Labels         []string   `json:"labels"`
		Assignees      []Assignee `json:"assignees"`
		Comments_Count *int       `json:"comments_count"`
		IssueUrl       string     `json:"issue_url"`
	}

	// get the token from the authentication header
	token := request_context.GetHeader("Authorization")

	//setup the github api
	github_api := setupGithubApi(token)

	// get the repository from the query string
	full_repository := request_context.Query("full_repository")

	// split the repository into the username and repository name
	repository_split := strings.Split(full_repository, "/")

	// confirm the repository is in the format username/repository
	if !isRepositoryFormat(repository_split) {
		request_context.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "repository must be in the format username/repository"})
		return
	}

	username := repository_split[0]
	repository := repository_split[1]

	// get the issues for the repository
	issues := getIssuesForRepository(username, repository, github_api)

	if issues == nil {
		request_context.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "error getting issues for repository"})
		return
	}

	// create a slice to hold the issue json
	var issues_json []Issue

	// loop through the issues
	for _, issue := range issues {

		// check the node ID is not PR_ prefix, if it is, skip the item
		if strings.HasPrefix(issue.GetNodeID(), "PR_") {
			continue
		}

		// get the array of labels
		labels := issue.Labels

		// create a slice to hold the labels
		var labels_array []string

		// loop through the labels
		for _, label := range labels {

			// append the label name to the slice
			labels_array = append(labels_array, label.GetName())
		}

		// get the array of assignees
		assignees := issue.Assignees

		// create a slice to hold the assignees
		var assignees_array []Assignee

		// loop through the assignees
		for _, assignee := range assignees {

			// get the assignee name
			assignee_name := assignee.GetLogin()

			// get the assignee url
			assignee_url := assignee.GetHTMLURL()

			// get the assignee avatar url
			assignee_avatar_url := assignee.GetAvatarURL()

			// create an assignee structure
			assignee_structure := Assignee{
				Login:      assignee_name,
				Html_Url:   assignee_url,
				Avatar_Url: assignee_avatar_url,
			}

			// append the assignee to the slice
			assignees_array = append(assignees_array, assignee_structure)

		}

		formatted_issue := Issue{
			// get the issue title
			Title: issue.GetTitle(),
			// get the issue body
			Body: issue.GetBody(),
			// get the issue number
			Number: issue.GetNumber(),
			// get the issue state
			State: issue.GetState(),
			// get the issue labels
			Labels: labels_array,
			// get the issue assignees
			Assignees: assignees_array,
			// get the issue comments count
			Comments_Count: issue.Comments,
			// get the issue url
			IssueUrl: issue.GetHTMLURL(),
		}

		// append the issue to the slice
		issues_json = append(issues_json, formatted_issue)
	}

	// return the issues in a json response
	request_context.JSON(http.StatusOK, gin.H{
		"total_count": len(issues_json),
		"data":        []any{issues_json}})

}

// function to confirm the repository is in the format username/repository
func isRepositoryFormat(repository_split []string) bool {

	// return true if the repository is in the format username/repository
	return len(repository_split) == 2
}

// function to get the issues for the repository
func getIssuesForRepository(
	username string,
	repository string,
	github_api *github.Client) []*github.Issue {

	// create a slice to hold all the issues
	var all_issues []*github.Issue

	// create an options structure for the issues request
	opt := &github.IssueListByRepoOptions{
		// set the state to all
		State: "open",
		// set the per page to 10
		ListOptions: github.ListOptions{PerPage: 100},
		// set the sort to created
		Sort: "created",
		// set the direction to desc
		Direction: "desc",
	}

	for {

		// get the issues for the repository
		issues, resp, err := github_api.Issues.ListByRepo(CTX, username, repository, opt)

		// confirm there was no error
		if err != nil {
			return nil
		}

		// append the issues to the slice
		all_issues = append(all_issues, issues...)

		// if there are no more pages, break out of the loop
		if resp.NextPage == 0 {
			break
		}

		// set the page to the next page
		opt.Page = resp.NextPage

	}

	// create an options structure for the issues request
	opt = &github.IssueListByRepoOptions{
		// set the state to all
		State: "closed",
		// set the per page to 10
		ListOptions: github.ListOptions{PerPage: 100},
		// set the sort to created
		Sort: "created",
		// set the direction to desc
		Direction: "desc",
	}

	for {

		// get the issues for the repository
		issues, resp, err := github_api.Issues.ListByRepo(CTX, username, repository, opt)

		// confirm there was no error
		if err != nil {
			return nil
		}

		// append the issues to the slice
		all_issues = append(all_issues, issues...)

		// if there are no more pages, break out of the loop
		if resp.NextPage == 0 {
			break
		}

		// set the page to the next page
		opt.Page = resp.NextPage

	}

	// return the issues
	return all_issues

}
