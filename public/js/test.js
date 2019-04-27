// Tests
function testChat() {
    if (typeof (Worker) !== "undefined") {
        if (typeof (w) == "undefined") {
            w = new Worker("/static/js/simulate.js");
        }
        w.onmessage = function (event) {
            msg.value = event.data;
            send();
        };
    } else {
        appendLog("<p>Sorry! No Web Worker support.</p>");
    }
}

// For debugging
function runTests() {
    var create = function () {
        console.log("Creating rooms")
        createRoom("default chat", "gay fuckfest", "public", "");
        createRoom("title 2 private", "gay fuckfest", "private", "");
    }
    var retrieve = function () {
        console.log("Retrieve Rooms...");
        retrieveRoom("Default Chat");
        retrieveRoom("title 2 private");
    }

    var update = function () {
        console.log("update rooms");
        putRoom("Default Chat", "changed to private", "private", "secret");
        retrieveRoom("Default Chat");
    }

    var del = function () {
        console.log("deleting rooms...")
        deleteRoom("Default Chat")
        deleteRoom("title 2 private")
    }

    // call 
    create()
    retrieve()
    update()
    del()
    alert("Completed tests successfully!")
}