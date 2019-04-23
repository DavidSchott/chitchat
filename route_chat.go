package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/DavidSchott/chitchat/data"
)

func sseActionHandler(w http.ResponseWriter, r *http.Request) {
	// read in request
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	// create ChatEvent obj
	var ce data.ChatEvent
	json.Unmarshal(body, &ce)

	// Fetch time
	ce.Timestamp = time.Now()

	// Perform requested action
	switch ce.EventType {
	case data.Unsubscribe:
		unsubscribe(w, r, &ce)
	case data.Subscribe:
		subscribe(w, r, &ce)
	default:
		broadcast(w, r, &ce)
	}
}

func broadcast(w http.ResponseWriter, r *http.Request, c *data.ChatEvent) {
	if cr, err := data.RetrieveID(c.RoomID); err == nil {
		flusher, _ := w.(http.Flusher)
		cr.Broker.Notifier <- formatEventData(c.Msg, c.User, c.Color)
		flusher.Flush()
	}
}

func subscribe(w http.ResponseWriter, r *http.Request, c *data.ChatEvent) {
	if cr, err := data.RetrieveID(c.RoomID); err == nil {
		// Add client
		client := &data.Client{
			Username: c.User,
			Color:    c.Color,
		}
		cr.Clients[c.User] = client
		info("Adding client to Chatroom: ", cr.Clients[c.User])
		cr.Broker.Notifier <- formatEventData(fmt.Sprintf("%s entered the room.", c.User), c.User, c.Color)
	}
}

func unsubscribe(w http.ResponseWriter, r *http.Request, c *data.ChatEvent) {
	if cr, err := data.RetrieveID(c.RoomID); err == nil {
		flusher, _ := w.(http.Flusher)
		cr.Broker.Notifier <- formatEventData(fmt.Sprintf("%s left the room.", c.User), c.User, c.Color)
		// Remove Client from tracked list
		delete(cr.Clients, c.User)
		info(fmt.Sprintf("Unsubscribing %s in room %d", c.User, cr.ID))
		flusher.Flush()
	}
}

// Upgrade to a sse connection
// Add to active chat session
func sseHandler(w http.ResponseWriter, r *http.Request) {
	if id, err := strconv.Atoi(path.Base(r.URL.Path)); err != nil {
		warning("Error creating sse for ", id, " Reason: ", err)
	} else {
		if cr, err := data.RetrieveID(id); err == nil {
			// Do stuff
			// Make sure that the writer supports flushing.
			//
			flusher, _ := w.(http.Flusher)

			// Each connection registers its own message channel with the Broker's connections registry
			messageChan := make(chan []byte)

			// Signal the broker that we have a new connection
			cr.Broker.NewClients <- messageChan

			// Remove this client from the map of connected clients
			// when this handler exits.
			defer func() {
				cr.Broker.ClosingClients <- messageChan
			}()

			// Listen to connection close and un-register messageChan
			notify := w.(http.CloseNotifier).CloseNotify()

			for {
				select {
				case <-notify:
					info("Closed connection in ", cr.Title)
					return
				default:
					// Write to the ResponseWriter
					// Server Sent Events compatible
					fmt.Fprintf(w, "%s", <-messageChan)

					// Flush the data immediatly instead of buffering it for later.
					flusher.Flush()
				}
			}
		} else {
			error_message(w, r, "Critical error creating SSE: "+err.Error())
			danger("error creating SSE: ", err)
		}
	}
}
