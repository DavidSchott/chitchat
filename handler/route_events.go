package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DavidSchott/chitchat/data"
	"github.com/gorilla/mux"
)

// /chats/{titleOrID}/sse/broadcast
func sseActionHandler(w http.ResponseWriter, r *http.Request) (err error) {
	queries := mux.Vars(r)
	if titleOrID, ok := queries["titleOrID"]; ok {
		// Fetch room & authorize
		cr, err := data.CS.Retrieve(titleOrID)
		if err != nil {
			info("erroneous chats API request", r, err)
			return err
		}
		// read in request
		len := r.ContentLength
		body := make([]byte, len)
		r.Body.Read(body)
		// create ChatEvent obj
		var ce data.ChatEvent
		if err := json.Unmarshal(body, &ce); err != nil {
			warning("error parsing JSON chatevent", r.Body, err)
			return err
		}
		// Set timestamp and room ID
		ce.Timestamp = time.Now()
		ce.RoomID = cr.ID
		// Check for invalid/random input
		if ce.User == "" {
			return &data.APIError{
				Code:  304,
				Field: "username",
			}
		}
		/* Authorize
		if cr.Type != data.PublicRoom {
			// if isn't public room, authorize
			cookieSecret, err := r.Cookie("secret_cookie")
			if err != nil {
				warning("error attempting to authorize "+strconv.Itoa(cr.ID)+" by:", ce)
				return &data.APIError{
					Code:  304,
					Field: "password",
				}
			}
			if cookieSecret.Value != cr.Password {
				return &data.APIError{
					Code:  304,
					Field: "password",
				}
			}
		}*/

		// Perform requested action
		switch ce.EventType {
		case data.Unsubscribe:
			// Populate activity
			cr.Clients[ce.User].LastActivity = ce.Timestamp
			unsubscribe(w, r, &ce, cr)
		case data.Subscribe:
			subscribe(w, r, &ce, cr)
		default:
			// Populate activity
			cr.Clients[ce.User].LastActivity = ce.Timestamp
			broadcast(w, r, &ce, cr)
		}
	}
	return
}

func broadcast(w http.ResponseWriter, r *http.Request, c *data.ChatEvent, cr *data.ChatRoom) {
	flusher, _ := w.(http.Flusher)
	cr.Broker.Notifier <- formatEventData(c.Msg, c.User, c.Color)
	flusher.Flush()

}

func subscribe(w http.ResponseWriter, r *http.Request, c *data.ChatEvent, cr *data.ChatRoom) {
	// Add client
	client := &data.Client{
		Username:     c.User,
		Color:        c.Color,
		LastActivity: time.Now(),
	}
	if err := cr.AddClient(client); err != nil {
		warning("error adding client:", err.Error())
		ReportStatus(w, false, err.(*data.APIError))
		return
	}
	info("Adding client to Chatroom: ", c.User)
	go func() {
		time.Sleep(200 * time.Millisecond)
		cr.Broker.Notifier <- formatEventData(fmt.Sprintf("%s entered the room.", c.User), c.User, c.Color)
	}()
	return
}

func unsubscribe(w http.ResponseWriter, r *http.Request, c *data.ChatEvent, cr *data.ChatRoom) {
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

// Upgrade to a sse connection
// Add to active chat session
// GET /chats/{titleOrID}/sse/subscribe
func sseHandler(w http.ResponseWriter, r *http.Request) (err error) {
	queries := mux.Vars(r)
	if titleOrID, ok := queries["titleOrID"]; ok {
		// Fetch room & authorize
		cr, err := data.CS.Retrieve(titleOrID)
		if err != nil {
			info("erroneous chats API request", r, err)
			return err
		}
		/*if cr.Type != data.PublicRoom {
			// if isn't public room, authorize
			cookieSecret, err := r.Cookie("secret_cookie")
			if err != nil {
				warning("error attempting to authorize "+titleOrID+" by:", *r)
				return &data.APIError{
					Code:  304,
					Field: "secret",
				}
			}
			if cookieSecret.Value != cr.Password {
				return &data.APIError{
					Code:  304,
					Field: "secret",
				}
			}
		}*/
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
				return err
			default:
				// Write to the ResponseWriter
				// Server Sent Events compatible
				fmt.Fprintf(w, "%s", <-messageChan)

				// Flush the data immediatly instead of buffering it for later.
				flusher.Flush()
			}
		}
	} else {
		errorMessage(w, r, "Critical error creating SSE: "+err.Error())
		danger("error creating SSE: ", err)
	}

	return
}
