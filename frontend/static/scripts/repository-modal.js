// Function for getting handling interactions with the repository modal
function getRepositoryModal(){

}

// function validate the owner and repository provided by the user
function validateRepository() {

    if (isDebugging) {
        console.debug("Validating repository");
    }

    // get the updated values
    gitHubRepository = repoDisplay.value;
    gitHubOwner = ownerDisplay.value;

    // if the owner and repository are not empty, enable the button
    if (gitHubRepository != "" && gitHubOwner != "") {

        // enable the button
        repositoryModalUpdateButton.disabled = false;

    }
}




// map validateRepository to the text input fields for owner and repository
$(ownerDisplay).on("change", validateRepository);
$(repoDisplay).on("change", validateRepository);