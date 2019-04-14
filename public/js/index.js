$(document).ready(function () {
    $('#sidebarCollapse').on('click', function () {
        $('#sidebar').toggleClass('active');
        $(this).toggleClass('active');
    });
    console.log("ready");
});

function runTests() {
    console.log("Create Rooms...");
    createRoom("Title 1", 0, "david", "");
    createRoom("Title 2 Private", 1, "david_private", "!!123abc");
    console.log("Retrieve Rooms...");
    retrieveRoom("title 1");
    retrieveRoom("title 2 private");
    console.log("update rooms");
    putRoom("Title 1", 1, "new_user", "secret");
    retrieveRoom("title 1");
    deleteRoom("title 1")
    deleteRoom("title 2 private")
    alert("Completed tests successfully!")
}

// Globals
var user

// POST /chat/
function createRoom(title, classification, user, password) {
    $.post('/chat/', JSON.stringify({ title: title, name: user, classification: classification, password: password }), "json")
        .done(function (data) {
            console.log(data)
            if (!data.hasOwnProperty('error')) {
                console.log("successfully created room!")
            }
            else {
                displayAlert("Could not create room " + title);
            }
        })
        .fail(function (xhr) {
            reject("Error creating room " + title);
            console.log(xhr);
        });
}

// GET /chat/<id>
function retrieveRoom(title) {
    $.get('/chat/' + title)
        .done(function (data) {
            console.log(data)
            if (!data.hasOwnProperty('error')) {
                console.log("success")
                return data;
            }
            else {
                displayAlert("Could not retrieve chat room  " + title);
            }
        })
        .fail(function (xhr) {
            reject("Error fetching chat room " + title);
            console.log(xhr);
        });
}
// PUT /chat/<id>
function putRoom(title, classification, user, password) {
    $.ajax({
        url: "/chat/" + title,
        method: 'PUT',
        data: JSON.stringify({ title: title, name: user, classification: classification, password: password })
    })
        .done(function (data) {
            console.log(data)
            if (!data.hasOwnProperty('error')) {
                console.log("success")
                return data;
            }
            else {
                displayAlert("Could not update chat room  " + title);
            }
        })
        .fail(function (xhr) {
            reject("Error fetching chat room " + title);
            console.log(xhr);
        });
    }

    function deleteRoom(title) {
        $.ajax({
            url: "/chat/" + title,
            method: 'DELETE'
        })
            .done(function (data) {
                console.log(data)
                if (!data.hasOwnProperty('error')) {
                    console.log("successfully deleted " + title)
                    return data;
                }
                else {
                    displayAlert("Could not delete chat room  " + title);
                }
            })
            .fail(function (xhr) {
                reject("Error deleting chat room " + title);
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