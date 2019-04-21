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
func checkStreamingSupport(w http.ResponseWriter, r *http.Request) (supported bool) {
	return
}

func broadcast(w http.ResponseWriter, r *http.Request, c *data.ChatEvent) {
	if cr, err := data.RetrieveID(c.RoomID); err == nil {
		flusher, ok := w.(http.Flusher)

		// Check if streaming is supported
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//broker.Notifier <- []byte(fmt.Sprintf("Removed client. %d registered Clients", len(broker.Clients)))
		cr.Broker.Notifier <- []byte(c.Msg)
		flusher.Flush()
	}
}

func subscribe(w http.ResponseWriter, r *http.Request, c *data.ChatEvent) {
	if cr, err := data.RetrieveID(c.RoomID); err == nil {
		_, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Add client
		client := &data.Client{
			Username: c.User,
			Color:    c.Color,
		}
		cr.Clients[c.User] = client
		p(cr.Clients)
	}
}

func unsubscribe(w http.ResponseWriter, r *http.Request, c *data.ChatEvent) {
	if cr, err := data.RetrieveID(c.RoomID); err == nil {
		flusher, ok := w.(http.Flusher)

		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Remove Client from tracked list
		delete(cr.Clients, c.User)
		//messageChan := make(chan []byte)
		//cr.Broker.ClosingClients <- messageChan
		info(fmt.Sprintf("Unsubscribing %s in room %d", c.User, cr.ID))
		cr.Broker.Notifier <- []byte(fmt.Sprintf("%s left the room.", c.User))
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
			flusher, ok := w.(http.Flusher)

			if !ok {
				http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.Header().Set("Access-Control-Allow-Origin", "*")

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
					fmt.Fprintf(w, "data: %s\n\n", <-messageChan)

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
