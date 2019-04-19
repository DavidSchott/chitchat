var create = function () {
    console.log("Create Rooms...");
    createRoom("Default Chat", "public", "david", "");
    createRoom("Title 2 Private", "private", "david_private", "!!123abc");
    createRoom("forever public chat", "public", "alexa", "");
}

function append(){
    document.getElementById("chat-box").innerHTML += '<p>User 1: diLorem ipsum dolor sit amet, consectetur apisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliqu </p><p>User 5: msg 10 </p>';
}

// Create new room
function newRoom() {
    var title = document.getElementById("input-title").value;
    var classification = document.getElementById("input-type").value;
    var user = document.getElementById("input-user").value;
    var password = document.getElementById("input-password").value;
    var valid = validateForm(title, classification, user, password)
    if (valid) {
        // Submit new room
        createRoom(title, classification, user, password);
        // close modal
        $('#create-modal').modal('toggle');
        // display new chat room and join it
        return true;
    }
}

function setInnerContent(url,title=''){
    console.log(url+title)
    $.get(url+title)
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

function validateForm(title, classification, user, password) {
    // validate input
    return true;
}

$(document).ready(function () {
    $('#sidebarCollapse').on('click', function () {
        $('#sidebar').toggleClass('active');
        $(this).toggleClass('active');
    });
    console.log("ready");
});


// For debugging
function runTests() {
    var retrieve = function () {
        console.log("Retrieve Rooms...");
        retrieveRoom("Default Chat");
        retrieveRoom("title 2 private");
    }

    var update = function () {
        console.log("update rooms");
        putRoom("Default Chat", "private", "new_user", "secret");
        retrieveRoom("Default Chat");
    }

    var del = function () {
        console.log("deleting rooms...")
        deleteRoom("Default Chat")
        deleteRoom("title 2 private")
    }

    // call 
    retrieve()
    update()
    //    del()
    alert("Completed tests successfully!")
}

// REST API calls
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
            console.log("Error creating room " + title);
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
            console.log("Error fetching chat room " + title);
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
            console.log("Error fetching chat room " + title);
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
            console.log("Error deleting chat room " + title);
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