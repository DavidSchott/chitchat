var username
var password
var msg
var log
var stream
var direction

$(document).ready(function () {
    // Get user session vars
    username = (Math.random() + 1).toString(36).substring(7).toString();
    color = "purple"
    password = "plaintext"
    // Get Room ID
    ID = window.location.pathname.split("/").pop();
    stream;
    msg = document.getElementById("msg");
    log = document.getElementById("chat-box");
    direction = "right";
});

var chat = function () {
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
                console.log('Opened connection');
                // TODO: Use send to announce?
                // TODO: Remove Modal once joined
                // TODO: Display errors
            };

            // Received server notification (chat message)
            stream.onmessage = function (evt) {
                console.log(evt);
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
            event = JSON.stringify({ type: action, name: user, id: parseInt(room), color: col, msg: message })
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
function pushBalon(message, user, time, col = "") {
    var item = document.createElement("div");
    item.setAttribute("data-is", user + " - " + time); // TODO
    var text = document.createElement("a");

    if (direction == "right") {
        item.className = "balon1 p-2 m-0 position-relative"
        // set float
        text.className = "float-right";
        // Toggle direction
        direction = "left";
    }
    else if (direction == "left") {
        item.className = "balon2 p-2 m-0 position-relative"
        // set float
        text.className = "float-left";
        // Toggle direction
        direction = "right";
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

function popBalon(){
    balon = document.getElementsByClassName("balon1");
    log.removeChild(balon[0]);
}

function updateTemplateStyle(user, color) {
    popBalon();
    pushBalon("Hey there! What's up?", user, new Date().toLocaleTimeString(), color);
    direction = "right";
}