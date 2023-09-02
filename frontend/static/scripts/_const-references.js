function getDebug() {
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
        // Perform your action for the localhost environment
        isDebugging = true;
    } 

    // print debugging enabled
    if (isDebugging) {
        console.debug("Debugging enabled");
    }
}

//
// General
//

isDebugging = false;
// when document is ready, get the debug value from the backend
$(document).ready(getDebug);