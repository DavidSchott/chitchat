

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
if(typeof(EventSource) == "undefined") {
    var item = document.createElement("div");
    item.innerHTML = "<b>Sorry, your browser does not support Server-Sent Events! "</b>";
    appendLog(item);
}
else{





    function startSession(id) {
        stream = new EventSource("/chat/sse/" + id);
        console.log("established stream: ", stream)
    }

    function sendClientEvent(action,user,room,col=""){
        event = JSON.stringify({ type: action, name: user, id: parseInt(room), color: col })
        console.log(event);
        $.post('/chat/sse/event', event, "json")
        .done(function (data) {
            console.log(data)
            if (!data.hasOwnProperty('error')) {
                console.log("successfully left!")
            }
            else {
                displayAlert("Could not leave source event.");
            }
        })
        .fail(function (xhr) {
            console.log("Could not leave source event.");
            console.log(xhr);
        });
        stream.close();
    }

    // Send a msg
    var send = function () {
        if (!stream) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        // TODO: Implement how to send data
        console.log("implement to send" , msg.value);
        //stream.send(msg.value);
        msg.value = "";
        return false;
    };
        // If supported, create EventSource
        startSession(ID);
        window.addEventListener('beforeunload', function() {
            sendClientEvent("leave", "david", ID);
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
        stream.onopen = function() {
            console.log('Opened connection');
          };

        // Listen for messages
        stream.onmessage = function (evt) {
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

        // Connection closed
        stream.onclose = function(code, reason) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed. Reason: " + "reason" + "</b>";
            appendLog(item);
        };

        // Connection error
        stream.onerror = function (event) {
            console.log(event);
          };
        }