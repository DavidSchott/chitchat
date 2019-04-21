package main

import (
	"net/http"
	"path"
	"strconv"

	"github.com/DavidSchott/chitchat/chat"
	"github.com/DavidSchott/chitchat/data"
)

// Upgrade to a WebSocket connection
// Add to active chat session
func socketHandler(w http.ResponseWriter, r *http.Request) {
	p(r)
	if id, err := strconv.Atoi(path.Base(r.URL.Path)); err != nil {
		warning("Error creating socket for ", id, " Reason: ", err)
	} else {
		if cr, err := data.RetrieveID(id); err == nil {
			chat.ServeWs(cr.Session, w, r)
		} else {
			error_message(w, r, "Critical error creating socket: "+err.Error())
			danger("error creating socket: ", err)
		}
	}
}

/*
func chatHandler(w http.ResponseWriter, r *http.Request) {
	var ce data.ChatEvent
	if err := json.NewDecoder(r.Body).Decode(&ce); err != nil {
		warning("ERROR: ", err)
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}
	defer r.Body.Close()
	go writer(&ce)
}

func writer(ce *data.ChatEvent) {
	data.Transmission <- ce
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := data.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		danger(err)
	}
	// register client
	data.Clients[ws] = true
}
*/
