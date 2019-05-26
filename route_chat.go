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

// POST /chat/sse/login
func login(w http.ResponseWriter, r *http.Request) (err error) {
	// read in request
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var c data.ChatEvent
	json.Unmarshal(body, &c)

	if cr, err := data.CS.RetrieveID(c.RoomID); err == nil {
		if cr.Type == data.PublicRoom {
			// Ignore public room
			ReportSuccess(w, true, "")
		} else if c.Password == cr.Password {
			// Success! Set Password
			cookieSecret := http.Cookie{
				Name:     "secret_cookie",
				Value:    c.Password,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookieSecret)
			ReportSuccess(w, true, "")
		} else {
			// send unauthorized error. TODO: send ReportSuccess false?
			return &data.APIError{
				Code:  304,
				Field: "secret",
			}
		}
	}
	return
}

// /chat/sse/event
func sseActionHandler(w http.ResponseWriter, r *http.Request) {
	// read in request
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	// create ChatEvent obj
	var ce data.ChatEvent
	json.Unmarshal(body, &ce)
	// Set timestamp
	ce.Timestamp = time.Now()

	// Fetch room & authenticate
	if cr, err := data.CS.RetrieveID(ce.RoomID); err == nil {
		// Check for invalid/random input
		if ce.User == "" || ce.Password != cr.Password {
			ReportSuccess(w, false, "Invalid credentials.")
			return
		}
		// Authenticate
		if cr.Type != data.PublicRoom {
			// if isn't public room, authenticate
			cookieSecret, err := r.Cookie("secret_cookie")
			if err != nil {
				ReportSuccess(w, false, "Error verifying credentials.")
				warning("error attempting to authenticate "+strconv.Itoa(cr.ID)+" by:", ce)
				return
			}
			if cookieSecret.Value != cr.Password {
				ReportSuccess(w, false, "Invalid credentials.")
				return
			}
		}

		// Perform requested action
		switch ce.EventType {
		case data.Unsubscribe:
			// Populate activity
			cr.Clients[ce.User].LastActivity = ce.Timestamp
			unsubscribe(w, r, &ce)
		case data.Subscribe:
			subscribe(w, r, &ce)
		default:
			// Populate activity
			cr.Clients[ce.User].LastActivity = ce.Timestamp
			broadcast(w, r, &ce)
		}
	}
}

func broadcast(w http.ResponseWriter, r *http.Request, c *data.ChatEvent) {
	if cr, err := data.CS.RetrieveID(c.RoomID); err == nil {
		flusher, _ := w.(http.Flusher)
		cr.Broker.Notifier <- formatEventData(c.Msg, c.User, c.Color)
		flusher.Flush()
	}
}

func subscribe(w http.ResponseWriter, r *http.Request, c *data.ChatEvent) {
	if cr, err := data.CS.RetrieveID(c.RoomID); err == nil {
		// Add client
		client := &data.Client{
			Username:     c.User,
			Color:        c.Color,
			LastActivity: time.Now(),
		}
		if err := cr.AddClient(client); err != nil {
			warning(err.Error())
			ReportSuccess(w, false, err.Error())
			return
		}
		info("Adding client to Chatroom: ", c.User)
		go func() {
			time.Sleep(200 * time.Millisecond)
			cr.Broker.Notifier <- formatEventData(fmt.Sprintf("%s entered the room.", c.User), c.User, c.Color)
		}()
	}
	return
}

func unsubscribe(w http.ResponseWriter, r *http.Request, c *data.ChatEvent) {
	if cr, err := data.CS.RetrieveID(c.RoomID); err == nil {
		flusher, _ := w.(http.Flusher)
		// Remove Client from tracked list
		//delete(cr.Clients, c.User)
		cr.RemoveClient(c.User)
		info(fmt.Sprintf("Unsubscribing %s in room %d", c.User, cr.ID))
		go func() {
			time.Sleep(200 * time.Millisecond)
			cr.Broker.Notifier <- formatEventData(fmt.Sprintf("%s left the room.", c.User), c.User, c.Color)
		}()
		flusher.Flush()
	}
}

// Upgrade to a sse connection
// Add to active chat session
// /chat/sse
func sseHandler(w http.ResponseWriter, r *http.Request) {
	if id, err := strconv.Atoi(path.Base(r.URL.Path)); err != nil {
		warning("Error creating sse for ", id, " Reason: ", err)
	} else {
		if cr, err := data.CS.RetrieveID(id); err == nil {
			if cr.Type != data.PublicRoom {
				// if isn't public room, authenticate
				cookieSecret, err := r.Cookie("secret_cookie")
				if err != nil {
					warning("error attempting to authenticate "+strconv.Itoa(id)+" by:", *r)
					ReportSuccess(w, false, "Error verifying credentials.")
					return
				}
				if cookieSecret.Value != cr.Password {
					ReportSuccess(w, false, "Unauthorized credentials.")
					return
				}
			}
			// Do stuff here
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
			// &data.APIError{Code: 301}
			error_message(w, r, "Critical error creating SSE: "+err.Error())
			danger("error creating SSE: ", err)
		}
	}
}
