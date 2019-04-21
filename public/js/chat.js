
var initChat = function () {
    var ID = window.location.pathname.split("/").pop();
    var socket;
    var msg = document.getElementById("msg");
    var log = document.getElementById("chat-box");
    var direction = "right";

    // Chat-related functions
    function appendLog(item) {
        console.log(item)
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    function startSession(id) {
        socket = new WebSocket("ws://" + document.location.host + "/chat/ws/" + id);
        console.log("established websocket: ", socket)
    }

    // Send a msg
    var send = function () {
        if (!socket) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        socket.send(msg.value);
        msg.value = "";
        return false;
    };

    // check if websockets are not supported
    if (!window["WebSocket"]) {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets. You cannot use this application, sorry!</b>";
        document.body.appendChild(item);
    }
    else {
        // If supported, create web socket
        startSession(ID); // TODO: Pass ID here,
        // Close the connection when the window is closed
        window.addEventListener('beforeunload', function() {
            socket.close();
         });
         
        // Handle msg send events
        document.getElementById("send-btn").onclick = send;
        msg.onkeydown = function (evt) {
            if (event.which == 13 || event.keyCode == 13) {
                send();
                return false;
            }
            return true;
        }

        // Connection opened
        socket.addEventListener('open', function (event) {
            socket.send('I joined!');
        });

        // Listen for messages
        socket.addEventListener('message', function (evt) {
            console.log('Message from server ', evt.data);
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                item.setAttribute("data-is", "username - 15:20"); // TODO
                // TODO: Add colors
                var text = document.createElement("a");
                if (direction == "right") {
                    item.className = "balon1 p-2 m-0 position-relative"
                    // set text
                    text.className = "float-right";
                    text.innerText = messages[i];
                    // Make it a child
                    item.appendChild(text);
                    // Toggle direction
                    direction = "left";
                }
                else if (direction == "left") {
                    item.className = "balon2 p-2 m-0 position-relative"
                    // set text
                    text.className = "float-left sohbet2";
                    text.innerText = messages[i];
                    // Make it a child
                    item.appendChild(text);
                    // Toggle direction
                    direction = "right";
                }
                appendLog(item);
            }
        });

        // Connection closed
        socket.addEventListener('close', function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        });

        // Connection closed
        socket.addEventListener('error', function (evt) {
            console.log("error", evt);
        });
    }
}
initChat();