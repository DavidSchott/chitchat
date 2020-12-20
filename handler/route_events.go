package handler

import (
	"net/http"

	"github.com/DavidSchott/chitchat/data"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow connections from any origin.
		},
	}
)

// Upgrade to a ws connection
// Add to active chat session
// GET /chats/{titleOrID}/ws
func webSocketHandler(w http.ResponseWriter, r *http.Request) (err error) {
	queries := mux.Vars(r)
	if titleOrID, ok := queries["titleOrID"]; ok {
		// Fetch room & authorize
		cr, err := data.CS.Retrieve(titleOrID)
		if err != nil {
			Warning("Error retrieving room", r, err)
			return err
		}
		// Do stuff here:
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			errorMessage(w, r, "Critical error creating WebSocket: "+err.Error())
			Danger("error creating WebSocket: ", err)
			return &data.APIError{Code: 301}
		}
		client := &data.Client{Room: cr, Conn: wsConn, Send: make(chan []byte)}
		client.Room.Broker.OpenClient <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.WritePump()
		go client.ReadPump()
	} else {
		return &data.APIError{Code: 101}
	}

	return
}
