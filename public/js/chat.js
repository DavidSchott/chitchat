var username = ""
var password = ""
var msg
var log
var stream
var ID;

$(document).ready(function () {
    ID = window.location.pathname.split("/").pop();
    msg = document.getElementById("msg");
    log = document.getElementById("chat-box");
});

var chat = function () {
    msg = document.getElementById("msg");
    log = document.getElementById("chat-box");
    // Check if SSE isn't supported
    if (typeof (EventSource) == "undefined") {
        var item = document.createElement("div");
        item.innerHTML = "<b>Sorry, your browser does not support Server-Sent Events!" + "</b>";
        appendLog(item);
        // TODO: Use WebSockets or PolyFill instead.
    }
    else {
        // Start EventSource
        register(ID);
        // Defer close
        window.addEventListener('beforeunload', function () {
            stream.close();
            sendClientEvent("leave", username, ID, "", color);
        });

        // Handle msg send events
        document.getElementById("send-btn").onclick = send;
        // Press enter to send chat
        msg.onkeydown = function (evt) {
            if (evt.which == 13 || evt.keyCode == 13) {
                send();
                return false;
            }
            return true;
        }
        //testChat();

        // Functions for Event sources
        // Start event source for current Room ID
        function startSession(id) {
            stream = new EventSource("/chat/sse/" + id);
            sendClientEvent("join", username, ID, "", color);
            console.log("established EventSource stream for " + id);
            return stream;
        }
        function register(id) {
            stream = startSession(id)
            // Event Source is opened
            stream.onopen = function () {
                console.log('Entered session');
                // TODO: Use send to announce?
                // TODO: Remove Modal once joined
                // TODO: Display errors
            };

            // Received server notification (chat message)
            stream.onmessage = function (evt) {
                json = JSON.parse(evt.data);
                var message = json.msg;
                var usr = json.name;
                var color = json.color;
                pushBalon(message, usr, new Date().toLocaleTimeString(), color);
            };
            /* TODO: Implement
            stream.addEventListener('join', function (e) {
                var data = JSON.parse(e.data);
                console.log('User login:' + data.username);
            }, false);
            stream.addEventListener('leave', function (e) {
                var data = JSON.parse(e.data);
                console.log('User login:' + data.username);
            }, false);
            */

            // Server connection closed
            stream.onclose = function (code, reason) {
                var item = document.createElement("div");
                item.innerHTML = "<b>Connection closed. Reason: " + "reason" + "</b>";
                appendLog(item);
            };

            stream.addEventListener('error', function (event) {
                console.log("Streaming Error:", event);
                switch (event.target.readyState) {

                    case EventSource.CONNECTING:
                        console.log('Reconnecting...');
                        break;

                    case EventSource.CLOSED:
                        console.log('Connection failed, will try to re-register');
                        register(ID);
                        break;
                }

            }, false);
        }
        // Send notification to server
        function sendClientEvent(action, user, room, message = "", col = "") {
            event = JSON.stringify({ type: action, name: user, id: parseInt(room), color: col, msg: message, password: password })
            $.post('/chat/sse/event', event, "json")
                .done(function (data) {
                })
                .fail(function (xhr) {
                    console.log("Failed sending client event:", event);
                });
        }

        // Submit a chat message to broadcast
        function send() {
            if (!stream) {
                return false;
            }
            if (!msg.value) {
                return false;
            }
            sendClientEvent("send", username, ID, msg.value, color);
            msg.value = "";
            return false;
        };
    }
}
// Chat-related cosmetic functions
function appendLog(item) {
    var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
    }
}

// Populate chat box
function pushBalon(message, user, time, col = "", direction = "") {
    var item = document.createElement("div");
    item.setAttribute("data-is", user + " - " + time); // TODO
    var text = document.createElement("a");

    if (user == username || direction == "right") {
        item.className = "balon1 p-2 m-0 position-relative"
        // set float
        text.className = "float-right";
    }
    else {
        item.className = "balon2 p-2 m-0 position-relative"
        // set float
        text.className = "float-left";
    }
    // Set txt
    text.innerText = message;
    // Add colors
    if (col != "") {
        text = applyColor(text, col);
    }
    // Make it a child
    item.appendChild(text);
    appendLog(item);
}
function applyColor(elem, color) {
    switch (color) {
        case "blue":
            elem.setAttribute("style", "background: #42a5f5; color: #ffffff !important;")
            break;
        case "red":
            elem.setAttribute("style", "background: #DC143C; color: #ffffff !important;")
            break;
        case "green":
            elem.setAttribute("style", "background: #2E8B57; color: #ffffff !important;")
            break;
        case "gray":
            elem.setAttribute("style", "background: #f1f1f1; color: #000 !important")
            break;
        case "turquoise":
            elem.setAttribute("style", "background: #40E0D0; color: #000 !important;")
            break;
        case "indigo":
            elem.setAttribute("style", "background: #4B0082; color: #ffffff !important;")
            break;
        case "magenta":
            elem.setAttribute("style", "background:#8B008B ; color: #ffffff !important;")
            break;
        case "black":
            elem.setAttribute("style", "background: #000000; color: #ffffff !important;")
            break;
        case "yellow":
            elem.setAttribute("style", "background: #FFD700; color: #000 !important;")
            break;
        case "orange":
            elem.setAttribute("style", "background: #FF8C00; color: #000 !important;")
            break;
        case "purple":
            elem.setAttribute("style", "background: #7386D5; color: #ffffff !important;")
            break;
        default:
            elem.setAttribute("style", "background: #7386D5; color: #ffffff !important;")
            break;
    }
    return elem
}

function popBalon() {
    balon = document.getElementsByClassName("balon1");
    log.removeChild(balon[0]);
}

function updateTemplateStyle(user, color) {
    popBalon();
    if (user == "") {
        user = "You"
    }
    pushBalon("Hey there! What's up?", user, new Date().toLocaleTimeString(), color, "right");
}

function loadChat() {
    username = document.getElementById("input-user").value;
    password = document.getElementById("inputPassword").value;
    color = document.getElementById("color-select").value;

    new Promise(
        function (resolve, reject) {
            setInnerContent("/chat/box/", ID, resolve, reject);
        })
        .then(function (result) {
            if (result.outcome) {
                chat();
            }
        })
        .catch(
            // Log the rejection reason (Room is invalid')
            function (outcome) {
                console.log(outcome);
                displayAlert(outcome.reason);
            });
}

function userExists(user, roomID, resolve = console.log, reject = console.log) {
    reject("TODO: look for user " + user + " in room " + roomID);
    
}

function checkPassword(password, resolve = console.log, reject = console.log) {
    reject("TODO: verify password = " + password);
}

function validateChatEntrance() {
    // Read in form
    var form = document.getElementsByClassName("form-signin")[0];
    var userDOM = document.getElementById("input-user");
    var passwordDOM = document.getElementById("inputPassword");
    var colorDOM = document.getElementById("color-select");
    valid = true;
    // validate fields look OK

    // check color is selected
    if (colorDOM.value.length < 1) {
        valid = false;
        colorDOM.setCustomValidity("No color selected");
    } else {
        colorDOM.setCustomValidity("");
    }
    // check user isn't empty
    if (userDOM.value.length < 1) {
        valid = false;
        userDOM.setCustomValidity("No user selected");
    } else {
        userDOM.setCustomValidity("");
    }
    // Check for duplicated user(s)
    if (valid) {
        new Promise(
            function (resolve, reject) {
                // Validate form looks good
                userExists(userDOM.value, ID, resolve, reject);
            })  // Check password
            .then(function (outcome) {
                // Define new promise to retrieve room
                console.log(outcome);
                new Promise(
                    function (resolve, reject) {
                        checkPassword(passwordDOM.value, resolve, reject);
                    })
                    .then(function (pwd) {
                        // Success! All conditions passed
                        console.log("joining chat room");
                        //loadChat();
                    }).catch(
                        function (reason) {
                            console.log(reason);
                            //displayAlert(reason);
                        });
            })
            .catch(
                function (issue) {
                    console.log("duplicate user " + userDOM.value);
                    document.getElementById('user-invalid-feedback').innerText = "Username already taken!";
                    userDOM.setCustomValidity("user-taken");
                    form.classList.add('was-validated');
                    reject(issue);
                });
    }
    form.classList.add('was-validated');
}