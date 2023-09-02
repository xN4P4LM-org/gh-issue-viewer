// Dependencies: jQuery
// function validate the owner and repository provided by the user
function validateRepository() {

    if (isDebugging) {
        console.debug("Validating repository");
    }

    repositoryInput = $("#repository-input");
    ownerInput = $("#owner-input");

    // get the updated values
    gitHubRepository = repositoryInput.value;
    gitHubOwner = ownerInput.value;

    // if the owner and repository are not empty, check the repository exists
    if (repositoryInput != "" && ownerInput != "") {

        // perform an ajax request to check if the repository exists
        $.ajax({
            url: "/api/repository/validate",
            type: "GET",
            // send the repository and owner as owner=&repository=
            data: {
                "owner": ownerInput,
                "repository": repositoryInput
            },
            success: function(data){

                if (isDebugging) {
                    console.debug(data);
                }
                
                if (data.Repository.valid === true) {
                    // enable the button
                    //repositoryModalUpdateButton.disabled = false;
                }

            },
            error: function(){
                // disable the button
                //repositoryModalUpdateButton.disabled = true;
                
            }
        });

    }else {

        // disable the button
        //repositoryModalUpdateButton.disabled = true;

    }
}

// validate user
function validateUser() {

    if (isDebugging) {
        console.debug("Validating user");
    }

    // get the updated values
    ownerInput = $("#owner-input").value;

    // if the owner is not empty, check the user exists
    if (ownerInput != "") {

        // perform an ajax request to check if the user exists
        $.ajax({
            url: "/api/user/validate",
            type: "GET",
            data: {
                "user": ownerInput
            },
            content_type: "application/json",
        }).done(function (msg) {
            if (isDebugging) {
                console.debug(msg);
            }
            
            if (msg.user.valid === true) {
                // enable the button
                repositoryModalUpdateButton.disabled = false;
            }
        }).fail(function() {
            // disable the button
            repositoryModalUpdateButton.disabled = true;   
        });
    }
    else {
            
        // disable the button
        repositoryModalUpdateButton.disabled = true;
    
    }

}

// map validateRepository to the text input fields for owner and repository

// print debugging adding event listeners
if (isDebugging) {
    console.debug("Adding event listeners");
}

$("#owner-input").on("blur", validateUser);
$("#repository-input").on("blur", validateRepository);