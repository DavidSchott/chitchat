var room
var title = ""
var description = ""
var visibility = ""
var password = ""

$(document).ready(function () {
    $('#sidebarCollapse').on('click', function () {
        $('#sidebar').toggleClass('active');
        $(this).toggleClass('active');
    });
    console.log("ready");
});

var create = function () {
    console.log("Create Rooms...");
    createRoom("Default Chat", "Another 2nd default chat", "public", "david", "");
    createRoom("Title 2 Private", "This is a password-protected secret room!", "private", "david_private", "!!123abc");
    createRoom("forever public chat", "This chat will always be available to the public!", "public", "alexa", "");
    createRoom("Hidden Chat", "super top secret chat... Cool!", "hidden", "jeff", "uber-secret-password");
}
function titleExists(title, resolve, reject) {
    // Invert reject/resolve to check for duplicates
    retrieveRoom(title, reject, resolve);
}

function validateForm(resolve = console.log, reject = console.log) {
    // Read in form
    var form = document.getElementById("create-new-room");
    var titleDOM = document.getElementById("input-title");
    var visibilityDOM = document.getElementById("input-type");
    var passwordDOM = document.getElementById("input-password");
    valid = true;
    // validate fields look OK
    // Check password
    if (visibilityDOM.value != "public" && (passwordDOM.value.length < 8 || passwordDOM.value.length > 20)) {
        valid = false;
        passwordDOM.setCustomValidity("Invalid password");
    } else {
        passwordDOM.setCustomValidity("");
    }
    // Check title
    if (titleDOM.value.length < 1) {
        valid = false;
        document.getElementById('title-invalid-feedback').innerText = "Please provide a valid title!";
        titleDOM.setCustomValidity("Empty title");
    } else {
        titleDOM.setCustomValidity("");
    }
    if (valid) {
        new Promise(
            function (resolve, reject) {
                // Validate form looks good
                titleExists(titleDOM.value, resolve, reject);
            }
        ).then(function (exists) {
                titleDOM.setCustomValidity("");
                form.classList.add('was-validated');
                resolve(exists);
        }).catch(
            function (duplicate) {
                console.log("duplicate title " + titleDOM.value);
                document.getElementById('title-invalid-feedback').innerText = "Room already exists!";
                titleDOM.setCustomValidity("Duplicate title");
                form.classList.add('was-validated');
                reject(duplicate);
            });
    }else{
        form.classList.add('was-validated');
        reject(valid);
    }
    //password.setAttribute("invalid");
    // password.classList.add('invalid');
    // password.classList.add('is-invalid');
}

// Create new room
function newRoom() {
    title = document.getElementById("input-title").value;
    description = document.getElementById("input-description").value;
    visibility = document.getElementById("input-type").value;
    //var user = document.getElementById("input-user").value;
    password = document.getElementById("input-password").value;
    new Promise(
        function (resolve, reject) {
            // Validate form looks good
            validateForm(resolve, reject);
        }
    ).then(function (validated) {
        new Promise(
            function (resolve, reject) {
                // Submit new room
                createRoom(title, description, visibility, password, resolve, reject);
            }
        )
            // wait for confirmation that room is created
            .then(function (outcome) {
                // Define new promise to retrieve room
                console.log(outcome);
                new Promise(
                    function (resolve, reject) {
                        retrieveRoom(title, resolve, reject);
                    })
                    .then(function (room) {
                        // Success! created & retrieved room object
                        console.log("fetched room", room);
                        // redirect to chat room
                        window.location.href = "/chat/join/" + room.id
                    }).catch(
                        function (reason) {
                            console.log(reason);
                            //displayAlert(reason);
                        });
            }).catch(
                // Log the rejection reason (Room is invalid')
                function (reason) {
                    console.log(reason);
                    //displayAlert(reason);
                });
    }).catch(
        function (validationError) {
            // Form is not valid and can't create room
            console.log(validationError);
            //displayAlert(validationError);
        });

    // close modal, redirect
    //$('#create-modal').modal('toggle');
    return true;
}

function setInnerContent(url, id = '', resolve = console.log, reject = console.log) {
    $.get(url + id)
        .done(function (data) {
            if (!data.hasOwnProperty('error')) {
                document.getElementById("inner-content").innerHTML = data;
                resolve({ outcome: true });
            }
            else {
                //displayAlert("Could not retrieve chat room");
                reject({ outcome: false, reason: "Could not retrieve chat room " + id });
            }
        })
        .fail(function (xhr) {
            console.log("Error fetching chat room " + id, xhr);
            reject({ outcome: false, reason: "Error fetching chat room " + id });
        });
}
function updateTypeDescription(visibility) {
    var helpText = ""
    switch (visibility) {
        case "public":
            helpText = "Public rooms are fully open and available to join for everyone!"
            document.getElementById("password-option").hidden = true;
            document.getElementById("input-password").required = false;
            break
        case "private":
            helpText = "Private rooms require a <i>password</i> to enter."
            document.getElementById("password-option").hidden = false;
            document.getElementById("input-password").required = true;
            break
        case "hidden":
            helpText = "Secret rooms are <i>unlisted</i> private rooms."
            document.getElementById("password-option").hidden = false;
            document.getElementById("input-password").required = true;
            break
    }
    document.getElementById("typeHelpInline").innerHTML = helpText;
}

function displayAlert(msg) {
    $('#error-alert').html('<strong>' + msg + '</strong>');
    $('#error-alert').show();
    //$('.loading').hide();
}

function hideAlert() {
    $('#error-alert').hide();
}