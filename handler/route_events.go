package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/DavidSchott/chitchat/data"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// /chats/{titleOrID}/ws/broadcast
func wsEventHandler(w http.ResponseWriter, r *http.Request) (err error) {
	queries := mux.Vars(r)
	if titleOrID, ok := queries["titleOrID"]; ok {
		// Fetch room & authorize
		cr, err := data.CS.Retrieve(titleOrID)
		if err != nil {
			Info("erroneous chats API request", r, err)
			return err
		}
		// read in request
		len := r.ContentLength
		body := make([]byte, len)
		if _, err := r.Body.Read(body); err != nil {
			Danger("Error reading", r, err.Error())
		}
		// create ChatEvent obj
		var ce data.ChatEvent
		if err := json.Unmarshal(body, &ce); err != nil {
			Warning("error parsing JSON chatevent", r.Body, err)
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

		// Perform requested action
		switch ce.EventType {
		case data.Unsubscribe:
			// Populate activity
			cr.Clients[ce.User].LastActivity = ce.Timestamp
			unsubscribe(w, r, &ce, cr)
		case data.Subscribe:
			// Activity will be populated in subscribe
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
		Warning("error adding client:", err.Error())
		ReportStatus(w, false, err.(*data.APIError))
		return
	}
	Info("Adding client to Chatroom: ", c.User)
	go func() {
		time.Sleep(200 * time.Millisecond)
		cr.Broker.Notifier <- formatEventData(fmt.Sprintf("%s entered the room.", c.User), c.User, c.Color)
	}()
}

func unsubscribe(w http.ResponseWriter, r *http.Request, c *data.ChatEvent, cr *data.ChatRoom) {
	flusher, _ := w.(http.Flusher)
	// Remove Client from tracked list
	//delete(cr.Clients, c.User)
	if err := cr.RemoveClient(c.User); err != nil {
		Danger("Error removing client", err.Error())
	}
	Info(fmt.Sprintf("Unsubscribing %s in room %d", c.User, cr.ID))
	go func() {
		time.Sleep(200 * time.Millisecond)
		cr.Broker.Notifier <- formatEventData(fmt.Sprintf("%s left the room.", c.User), c.User, c.Color)
	}()
	flusher.Flush()
}

// Upgrade to a ws connection
// Add to active chat session
// GET /chats/{titleOrID}/ws/subscribe
func wsHandler(w http.ResponseWriter, r *http.Request) (err error) {
	queries := mux.Vars(r)
	if titleOrID, ok := queries["titleOrID"]; ok {
		// Fetch room & authorize
		cr, err := data.CS.Retrieve(titleOrID)
		if err != nil {
			Info("erroneous chats API request", r, err)
			return err
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
		done := r.Context().Done()

		for {
			select {
			case <-done:
				Info("Closed connection in ", cr.Title)
				return err
			default:
				// Write to the ResponseWriter
				// Server Sent Events compatible
				fmt.Fprintf(w, "%s", <-messageChan)

				// Flush the data immediately instead of buffering it for later.
				flusher.Flush()
			}
		}
	} else {
		errorMessage(w, r, "Critical error creating SSE: "+err.Error())
		Danger("error creating SSE: ", err)
	}

	return
}

// convenience function to be chained with another HandlerFunc
// Checks if streaming via Server-Side Events is supported by the device
func checkWebSocketSupport(h errHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			Danger("Unable to upgrade to websockets", err.Error())
			http.Error(w, "Unable to upgrade to websockets", http.StatusBadRequest)
			return
		}
		if err := h(w, r); err != nil {
			Warning("Error calling:", runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name())
		}
	}
}

func formatEventData(msg string, user string, color string) (data []byte) {
	json := strings.Join([]string{
		fmt.Sprintf("data: {\"msg\": \"%s\",", msg),
		fmt.Sprintf("\"name\": \"%s\",", user),
		fmt.Sprintf("\"color\": \"%s\"}\n", color),
		"\n\n",
	}, "")
	data = []byte(json)

	return
}
