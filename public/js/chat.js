$(document).ready(function () {


// Get user session vars
var username = "david"
var color = "black"
var password = "plaintext"
// Get Room ID
var ID = window.location.pathname.split("/").pop();
var stream;
var msg = document.getElementById("msg");
var log = document.getElementById("chat-box");
var direction = "right";

// Chat-related functions
function appendLog(item) {
    var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
    }
}
if (typeof (EventSource) == "undefined") {
    var item = document.createElement("div");
    item.innerHTML = "<b>Sorry, your browser does not support Server-Sent Events!" + "</b>";
    appendLog(item);
    // TODO: Use WebSockets or PolyFill instead.
}
else {
    // Functions for Event sources

    // Start event source for current Room ID
    function startSession(id) {
        stream = new EventSource("/chat/sse/" + id);
        console.log("established stream: ", stream);
    }
    // Send notification to server
    function sendClientEvent(action, user, room, message="", col = "") {
        event = JSON.stringify({ type: action, name: user, id: parseInt(room), color: col, msg:message })
        $.post('/chat/sse/event', event, "json")
            .done(function (data) {
                console.log(data)
            })
            .fail(function (xhr) {
                console.log("Could not leave source event.");
                console.log(xhr);
            });
    }

    // Send a chat message
    var send = function () {
        if (!stream) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        sendClientEvent("send",username,ID,msg.value,color);
        msg.value = "";
        return false;
    };

    // If supported, create EventSource
    startSession(ID);
    // Defer close
    window.addEventListener('beforeunload', function () {
        stream.close();
        sendClientEvent("leave", username, ID);
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

    // Event Source is opened
    stream.onopen = function () {
        sendClientEvent("join", username, ID, color)
        console.log('Opened connection');
        // TODO: Use send to announce?
        // TODO: Remove Modal once joined
        // TODO: Display errors
    };

    // Received server notification (chat message)
    stream.onmessage = function (evt) {
        console.log(evt);
        var message = evt.data;
        var item = document.createElement("div");
        item.setAttribute("data-is", "username - 15:20"); // TODO
        // TODO: Add colors
        var text = document.createElement("a");
        if (direction == "right") {
            item.className = "balon1 p-2 m-0 position-relative"
            // set text
            text.className = "float-right";
            text.innerText = message;
            // Make it a child
            item.appendChild(text);
            // Toggle direction
            direction = "left";
        }
        else if (direction == "left") {
            item.className = "balon2 p-2 m-0 position-relative"
            // set text
            text.className = "float-left sohbet2";
            text.innerText = message;
            // Make it a child
            item.appendChild(text);
            // Toggle direction
            direction = "right";
        }
        appendLog(item);
    };

    // Server connection closed
    stream.onclose = function (code, reason) {
        var item = document.createElement("div");
        item.innerHTML = "<b>Connection closed. Reason: " + "reason" + "</b>";
        appendLog(item);
    };

    // Connection error
    stream.onerror = function (event) {
        console.log("EventSource error", event);
    };
}
});