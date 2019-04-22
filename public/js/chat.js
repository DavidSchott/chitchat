$(document).ready(function () {
    // Get user session vars
    var username = (Math.random() + 1).toString(36).substring(7).toString();
    var color = "gray"
    var password = "plaintext"
    // Get Room ID
    var ID = window.location.pathname.split("/").pop();
    var stream;
    var msg = document.getElementById("msg");
    var log = document.getElementById("chat-box");
    var direction = "right";

    // Check if SSE isn't supported
    if (typeof (EventSource) == "undefined") {
        var item = document.createElement("div");
        item.innerHTML = "<b>Sorry, your browser does not support Server-Sent Events!" + "</b>";
        appendLog(item);
        // TODO: Use WebSockets or PolyFill instead.
    }
    else {
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
            console.log(message, usr, color);
            pushBalon(message, usr, new Date().toLocaleTimeString(), color);
        };
        // TODO: Implement
        stream.addEventListener('join', function (e) {
            var data = JSON.parse(e.data);
            console.log('User login:' + data.username);
        }, false);


        stream.addEventListener('leave', function (e) {
            var data = JSON.parse(e.data);
            console.log('User login:' + data.username);
        }, false);


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
                    console.log('Connection failed, will not reconnect');
                    break;
            }

        }, false);

        // Functions for Event sources
        // Start event source for current Room ID
        function startSession(id) {
            stream = new EventSource("/chat/sse/" + id);
            sendClientEvent("join", username, ID, "", color);
            console.log("established EventSource stream for " + id);
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
            sendClientEvent("send", username, ID, msg.value);
            msg.value = "";
            return false;
        };

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
            case "purple":
                elem.style.background = '#7386D5';
                elem.style.color = '#ffffff !important';
                break;
            case "blue":
                elem.style.background = '#42a5f5';
                elem.style.color = '#ffffff !important';
                break;
            case "red":
                elem.style.background = '#DC143C';
                elem.style.color = '#ffffff !important';
                break;
            case "green":
                elem.style.background = '#2E8B57';
                elem.style.color = '#ffffff !important';
                break;
            case "gray":
                elem.style.background = '#f1f1f1';
                elem.style.color = "#000 !important";
                break;
            case "turquoise":
                elem.style.background = '#40E0D0';
                elem.style.color = '#000 !important';
                break;
            case "indigo":
                elem.style.background = '#4B0082';
                elem.style.color = '#ffffff !important';
                break;
            case "magenta":
                elem.style.background = '#8B008B';
                elem.style.color = '#ffffff !important';
                break;
            case "black":
                elem.style.background = '#000000';
                elem.style.color = '#ffffff !important';
                break;
            case "yellow":
                elem.style.background = '#FFD700';
                elem.style.color = '#ffffff !important';
                break;
            case "orange":
                elem.style.background = '#FF8C00';
                elem.style.color = '#000 !important';
                break;
        }
        return elem
    }
});

