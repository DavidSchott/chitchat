// run random tests
function simulateChats() {
    setTimeout("simulateChats()", 500);
    postMessage((Math.random() + 1).toString(36).substring(7).toString());
}
simulateChats();