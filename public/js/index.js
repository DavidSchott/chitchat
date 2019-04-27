var room
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

// Create new room
function newRoom() {
    var title = document.getElementById("input-title").value;
    var description = document.getElementById("input-description").value;
    var classification = document.getElementById("input-type").value;
//    var user = document.getElementById("input-user").value;
    var password = document.getElementById("input-password").value;

    new Promise(
        function (resolve, reject) {
            // Submit new room
            createRoom(title, description, classification, password, resolve, reject);
            //checkRequest(title, description, classification, password, resolve, reject);
        }
    )
    // wait for resolution
    .then(function (outcome) {
        // Define new promise to retrieve room
        new Promise(
            function (resolve, reject) {
            retrieveRoom(title, resolve, reject);
        })
        .then(function (room) {
            // Success! Created & retrieved room object
            console.log("fetched room", room);
            window.location.href = "/chat/join/" + room.id
        }).catch(
            function (reason) {
                console.log(reason);
                displayAlert(reason);
            });
    })
        .catch(
            // Log the rejection reason (Room is invalid')
            function (reason) {
                console.log(reason);
                displayAlert(reason);
            });

    // close modal, redirect
    $('#create-modal').modal('toggle');
    return true;
}

function checkRequest(title, description, classification, password) {
    // validate input
    return true;
}

function setInnerContent(url, id = '') {
    $.get(url + id)
        .done(function (data) {
            console.log(data);
            if (!data.hasOwnProperty('error')) {
                document.getElementById("inner-content").innerHTML = data;
            }
            else {
                displayAlert("Could not retrieve chat room");
            }
        })
        .fail(function (xhr) {
            console.log("Error fetching chat room list");
            console.log(xhr);
        });
}



function displayAlert(msg) {
    $('#error-alert').html('<strong>' + msg + '</strong>');
    $('#error-alert').show();
    //$('.loading').hide();
}

function hideAlert() {
    $('#error-alert').hide();
}