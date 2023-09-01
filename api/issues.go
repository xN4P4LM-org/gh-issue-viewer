package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v54/github"
)

//
// Issues
//

type RequestOverview struct {
	Total_Count    uint64 `json:"total_count"`
	PR_Count       uint64 `json:"pr_count"`
	Issue_Count    uint64 `json:"issue_count"`
	Provided_Limit bool   `json:"provided_limit"`
	Limit_Count    uint64 `json:"limit_count,omitempty"`
}

// function to get all issues for a given repository
func getIssues(request_context *gin.Context) {

	// define the assignee structure
	type Assignee struct {
		Login      string `json:"login"`
		Html_Url   string `json:"html_url"`
		Avatar_Url string `json:"avatar_url"`
	}

	type Labels struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
	}

	// strut to hold the branch information
	type Branch struct {
		Ref_Name       string `json:"ref_name"`
		Repo_Name      string `json:"repo_name"`
		Repo_Full_Name string `json:"repo_full_name"`
		Is_Fork        bool   `json:"is_fork"`
		Is_Default     bool   `json:"is_default"`
		Link_To_Branch string `json:"link_to_branch"`
	}

	type PR struct {
		Title       string   `json:"title"`
		Number      int      `json:"number"`
		Html_url    string   `json:"html_url"`
		State       string   `json:"state"`
		Body        string   `json:"body"`
		Labels      []Labels `json:"labels,omitempty"`
		Base_Branch Branch   `json:"base_branch"`
		Head_Branch Branch   `json:"head_branch"`
	}

	// define the issue structure
	type Issue struct {
		Title          string     `json:"title"`
		Body           string     `json:"body"`
		Number         int        `json:"number"`
		State          string     `json:"state"`
		Labels         []Labels   `json:"labels,omitempty"`
		Assignees      []Assignee `json:"assignees,omitempty"`
		Comments_Count *int       `json:"comments_count"`
		IssueUrl       string     `json:"issue_url"`
		Associated_PR  PR         `json:"associated_pull_request,omitempty"`
	}

	//setup the github api
	github_api, err := setupGithubApi()

	// confirm there was no error
	if err != nil {
		request_context.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "error setting up github api"})
		return
	}

	owner := request_context.Query("owner")
	repository := request_context.Query("repository")
	request_limit, err := strconv.ParseUint((request_context.Query("limit")), 10, 32)

	// confirm there was no error
	if err != nil {
		// if error is due to empty string, set the request limit to 0
		if strings.Contains(err.Error(), "strconv.ParseUint: parsing \"\": invalid syntax") {
			request_limit = 0

		} else {
			request_context.JSON(
				http.StatusBadRequest,
				gin.H{
					"msg":   "error parsing limit, must be a positive integer",
					"error": err.Error()})
			return
		}
	}

	// get the issues for the repository
	issues, headers, request_overview, err := getIssuesForRepository(owner,
		repository,
		github_api,
		request_limit)

	// confirm there was no error
	if err != nil {
		request_context.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error(), "headers": headers})
		return
	}

	// create a slice to hold the issue json
	var issues_json []Issue

	// loop through the issues
	for _, issue := range issues {

		// check the node ID is not PR_ prefix, if it is, skip the item
		if strings.HasPrefix(issue.GetNodeID(), "PR_") || strings.Contains(issue.GetHTMLURL(), "/pull/") {
			continue
		}

		var labels_array []Labels
		var assignees_array []Assignee
		var pr_struct PR

		// loop through the labels to create the label strut
		for _, label := range issue.Labels {

			// append the label name to the slice
			labels_array = append(labels_array, Labels{
				Name:        label.GetName(),
				Description: label.GetDescription(),
				Color:       label.GetColor(),
			})
		}

		// loop through the assignees to create the assignee struct
		for _, assignee := range issue.Assignees {

			// append the assignee to the slice
			assignees_array = append(assignees_array, Assignee{
				Login:      assignee.GetLogin(),
				Html_Url:   assignee.GetHTMLURL(),
				Avatar_Url: assignee.GetAvatarURL(),
			})

		}

		// if the associated pull request is not an empty string, get the pull request
		if issue.GetPullRequestLinks().GetURL() != "" {

			// get the pull request url
			pr_url := issue.GetPullRequestLinks().GetURL()

			// get the pull request number, owner, and repository from the url
			pr_owner := strings.Split(pr_url, "/")[4]
			pr_repository := strings.Split(pr_url, "/")[5]
			pr_number, err_convert := strconv.ParseInt(strings.Split(pr_url, "/")[6], 10, 32)

			// confirm there was no error
			if err_convert != nil {
				request_context.JSON(
					http.StatusInternalServerError,
					gin.H{"error": err_convert.Error(), "headers": headers})
				return
			}

			// get the pull request
			pr, _, err := github_api.PullRequests.Get(CTX, pr_owner, pr_repository, int(pr_number))

			// confirm there was no error
			if err != nil {
				request_context.JSON(
					http.StatusInternalServerError,
					gin.H{"error": err.Error(), "headers": headers})
				return
			}

			// create a slice to hold the pr labels
			var pr_labels_array []Labels

			// loop through the pr labels to create the label strut
			for _, pr_label := range pr.Labels {
				// append the label name to the slice
				pr_labels_array = append(pr_labels_array, Labels{
					Name:        pr_label.GetName(),
					Description: pr_label.GetDescription(),
					Color:       pr_label.GetColor(),
				})
			}

			var head_branch Branch
			var base_branch Branch

			// get the head branch information
			head_branch.Ref_Name = pr.Head.GetRef()
			head_branch.Repo_Name = pr.Head.GetRepo().GetName()
			head_branch.Repo_Full_Name = pr.Head.GetRepo().GetFullName()
			head_branch.Is_Fork = pr.Head.GetRepo().GetFork()
			if pr.Head.GetRepo().GetDefaultBranch() == pr.Head.GetRef() {
				head_branch.Is_Default = true
			} else {
				head_branch.Is_Default = false
			}
			head_branch.Link_To_Branch = pr.Head.GetRepo().GetHTMLURL() + "/tree/" + pr.Head.GetRef()

			// get the base branch information
			base_branch.Ref_Name = pr.Base.GetRef()
			base_branch.Repo_Name = pr.Base.GetRepo().GetName()
			base_branch.Repo_Full_Name = pr.Base.GetRepo().GetFullName()
			base_branch.Is_Fork = pr.Base.GetRepo().GetFork()
			if pr.Base.GetRepo().GetDefaultBranch() == pr.Base.GetRef() {
				base_branch.Is_Default = true
			} else {
				base_branch.Is_Default = false
			}
			base_branch.Link_To_Branch = pr.Base.GetRepo().GetHTMLURL() + "/tree/" + pr.Base.GetRef()

			// create the pr struct
			pr_struct = PR{
				Title:       pr.GetTitle(),
				Number:      pr.GetNumber(),
				Html_url:    pr.GetHTMLURL(),
				State:       pr.GetState(),
				Body:        pr.GetBody(),
				Labels:      pr_labels_array,
				Base_Branch: base_branch,
				Head_Branch: head_branch,
			}

		}

		formatted_issue := Issue{
			Title:          issue.GetTitle(),
			Body:           issue.GetBody(),
			Number:         issue.GetNumber(),
			State:          issue.GetState(),
			Labels:         labels_array,
			Assignees:      assignees_array,
			Comments_Count: issue.Comments,
			IssueUrl:       issue.GetHTMLURL(),
		}

		if issue.GetPullRequestLinks() != nil {
			// Add the associated pull request to the struct
			formatted_issue.Associated_PR = pr_struct
		} else {
			formatted_issue.Associated_PR = PR{}
		}

		// append the issue to the slice
		issues_json = append(issues_json, formatted_issue)
	}

	// return the issues in a json response
	request_context.JSON(http.StatusOK, gin.H{
		"request_overview": request_overview,
		"headers":          headers,
		"data":             []any{issues_json}})

}

// function to get the issues for the repository
func getIssuesForRepository(username string,
	repository string,
	github_api *github.Client,
	request_limit uint64) ([]*github.Issue, Headers, RequestOverview, error) {

	// create a slice to hold all the issues
	var all_issues []*github.Issue

	// create a request overview structure
	request_overview := RequestOverview{}

	var per_page int

	// if request limit is less than 100, set the per page to the request limit otherwise set it to 100
	if request_limit < 100 {
		per_page = int(request_limit)
	} else {
		per_page = 100
	}

	// create an options structure for the issues request
	opt := &github.IssueListByRepoOptions{
		// set the state to all
		State: "all",
		// set the per page to either the request limit or 100
		ListOptions: github.ListOptions{PerPage: per_page},
		// set the sort to created
		Sort: "created",
		// set the direction to desc
		Direction: "desc",
	}

	var issue_count uint64 = 0
	var pr_count uint64 = 0

	headers := Headers{}

	for {

		var total_count uint64

		// get the issues for the repository
		issues, resp, err := github_api.Issues.ListByRepo(CTX, username, repository, opt)

		// loop through the issues and count those that don't have a PR_ prefix
		for _, issue := range issues {

			// check the node ID is not PR_ prefix, if it is, skip the item
			if strings.HasPrefix(issue.GetNodeID(), "PR_") || strings.Contains(issue.GetHTMLURL(), "/pull/") {
				pr_count++
				continue
			}
			issue_count++
		}

		total_count = issue_count + pr_count

		// capture the headers
		headers.RateLimitRemaining = resp.Rate.Remaining
		headers.RequestID = resp.Header.Get("X-GitHub-Request-Id")
		headers.Date = resp.Header.Get("Date")

		// confirm there was no error
		if err != nil {

			request_overview.Total_Count = total_count
			request_overview.PR_Count = pr_count
			request_overview.Issue_Count = issue_count
			request_overview.Provided_Limit = request_limit != 0
			request_overview.Limit_Count = request_limit

			return nil, headers, request_overview, err
		}

		// append the issues to the slice
		all_issues = append(all_issues, issues...)

		// if there are no more pages, break out of the loop
		if resp.NextPage == 0 {
			break
		}

		// if the request limit is set and the issue count is greater than or equal to the request limit, break out of the loop
		if request_limit != 0 && total_count >= request_limit {
			break
		}

		// set the page to the next page
		opt.Page = resp.NextPage

	}

	request_overview.Total_Count = issue_count + pr_count
	request_overview.PR_Count = pr_count
	request_overview.Issue_Count = issue_count
	request_overview.Provided_Limit = request_limit != 0
	if request_limit != 0 {
		request_overview.Limit_Count = request_limit
	}

	// return the issues
	return all_issues,
		headers,
		request_overview,
		nil

}
