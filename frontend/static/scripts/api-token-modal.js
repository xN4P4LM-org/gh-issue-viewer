    // create a function to send the token to the backend and receive a response
    function getUsername() { 

        // get the token from the input box with element id GitHub-token
        token = $("#GitHub-token").val();

        fetch("/api/user", {
            method: "GET",
            headers: {
                "Authorization": `${token}`,
                "Content-Type": "application/json"
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            // update the username in the span with element id GitHub-username
            $("#GitHub-username").html(data["username"]);

            // write the login to console
            console.log(data);
            
            // change the api-token-modal-button to say "Change Token"
            $("#api-token-modal-pop-up-button").html("Change Token");

            // change the button to the color green and remove all other colors
            $("#api-token-modal-pop-up-button").removeClass("btn-warning");
            $("#api-token-modal-pop-up-button").removeClass("btn-danger");
            $("#api-token-modal-pop-up-button").addClass("btn-success");

            // change the button to say "Change Token" and remove all other colors
            $("#api-token-modal-update-token").html("Change Token");
            $("#api-token-modal-update-token").removeClass("btn-primary");
            $("#api-token-modal-update-token").removeClass("btn-danger");
            $("#api-token-modal-update-token").addClass("btn-success");

            // close the modal after the token is entered
            $("#api-token-modal").modal("hide");

        })
        .catch(error => {
            console.log('There was a problem with the fetch operation:', error.message);
            
            // set the token button to say "Invalid Token" and change the color to red
            $("#api-token-modal-update-token").html("Invalid Token");
            $("#api-token-modal-update-token").removeClass("btn-primary");
            $("#api-token-modal-update-token").addClass("btn-danger");

            $("#api-token-modal-pop-up-button").removeClass("btn-warning");
            $("#api-token-modal-pop-up-button").removeClass("btn-danger");
            $("#api-token-modal-pop-up-button").addClass("btn-danger");
        });
    }

    // map the function to the button
    $("#api-token-modal-update-token").on("click", getUsername);