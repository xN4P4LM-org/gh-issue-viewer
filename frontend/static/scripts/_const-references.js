function getDebug() {
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
        // Perform your action for the localhost environment
        isDebugging = true;
    } 
}

//
// General
//

isDebugging = false;
// when document is ready, get the debug value from the backend
$(document).ready(getDebug);

//
// API Token Modal
//

const apiTokenModal = document.getElementById("api-token-modal");
const apiModalButton = document.getElementById("api-token-modal-pop-up-button");
const apiSaveButton = document.getElementById("api-token-modal-update-token");
const apiInputField = document.getElementById("GitHub-token");
apiToken = apiInputField.value;
const apiTokenModalInput = document.getElementById("api-token-modal-input");
const publicReposOnlyButton = document.getElementById("public-repos-only-button");
const publicReposOnlyTooltip = document.getElementById("public-repos-only-tooltip");
const credentialsTooltip = document.getElementById("credentials-tooltip");
const credentialsTooltipDefault = $("#credentials-tooltip").html();

usePat = false;

//
// Repository Modal
//

const repositoryModalLabel = document.getElementById("repository-modal-label");
const ownerInput = document.getElementById("owner-input");
const repositoryInput = document.getElementById("repository-input");
const repositoryModalButtonLabel = document.getElementById("repository-modal-button-label");
const repositoryModalUpdateButton = document.getElementById("repository-modal-update-button");


//
// Username container
//

const gitHubUsername = document.getElementById("GitHub-username");
const usernameDisplay = document.getElementById("username-display");
const usernameButtonTooltip = document.getElementById("username-button-tooltip");

//
// Repository container
//

var gitHubRepository
var gitHubOwner
const repositoryContainer = document.getElementById("repository-container");
const repositoryDisplay = document.getElementById("repository");
const repoDisplay = document.getElementById("repository-display");
const ownerDisplay = document.getElementById("owner-display");

//
// issues container
//

const issuesContainer = document.getElementById("issues-container");
const issueCount = document.getElementById("issue-count");
const issueSort = document.getElementById("issue-sort");
const issueFilter = document.getElementById("issue-filter");
const issueFilterValue = document.getElementById("issue-filter-value");
const issueResults = document.getElementById("issue-results");
