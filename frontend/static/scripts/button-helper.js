// update button css
function updateButton(id,css_class,text){

    // create a variable to hold all possible classes for a given button
    var classes = 
    ["btn-primary", 
    "btn-secondary",
    "btn-success",
    "btn-danger",
    "btn-warning",
    "btn-info",
    "btn-light",
    "btn-dark",
    "btn-link"];
    
    // remove all classes from the element with the given id
    $(id).removeClass(classes.join(" "));

    // add the given class to the element with the given id
    $(id).addClass(css_class);
    
    // set the text of the element with the given id
    $(id).html(text);
}