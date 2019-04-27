// REST API calls
// POST /chat/
 function createRoom(title, description, classification, password, resolve, reject) {
    $.post('/chat/', JSON.stringify({ title: title, description: description, classification: classification, password: password }), "json")
        .done(function (data) {
            if (!data.hasOwnProperty('error')) {
                console.log("successfully created room!")
                resolve(data)
            }
            else {
                reject("Could not create room " + title)
            }
        })
        .fail(function (xhr) {
            reject("Could not create room " + title)
        });
}

// GET /chat/<id>
function retrieveRoom(title,resolve,reject) {
    $.get('/chat/' + title)
        .done(function (data) {
            if (!data.hasOwnProperty('error')) {
                resolve(data);
            }
            else {
                reject("Could not retrieve chat room  " + title)
            }
        })
        .fail(function (xhr) {
            reject("Error fetching chat room  " + title)
        });
}

// GET /chat/<id>
function retrieveRoomID(ID) {
    $.get('/chat/' + ID)
        .done(function (data) {
            console.log(data)
            if (!data.hasOwnProperty('error')) {
                console.log("success")
                return data;
            }
            else {
                displayAlert("Could not retrieve chat room  " + ID);
            }
        })
        .fail(function (xhr) {
            console.log("Error fetching chat room " + ID);
            console.log(xhr);
        });
}
// PUT /chat/<id>
function putRoom(title, description, classification, password) {
    $.ajax({
        url: "/chat/" + title,
        method: 'PUT',
        data: JSON.stringify({ title: title, description: description, classification: classification, password: password })
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

function deleteRoomID(ID) {
    $.ajax({
        url: "/chat/" + ID,
        method: 'DELETE'
    })
        .done(function (data) {
            console.log(data)
            if (!data.hasOwnProperty('error')) {
                console.log("successfully deleted " + ID)
                return data;
            }
            else {
                displayAlert("Could not delete chat room  " + ID);
            }
        })
        .fail(function (xhr) {
            console.log("Error deleting chat room " + ID);
            console.log(xhr);
        });
}