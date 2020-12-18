package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client represents a user in a ChatRoom
type Client struct {
	Username     string    `json:"username"`
	Color        string    `json:"color"` // Not being used due to possible connection failure...
	LastActivity time.Time `json:"last_activity"`
	// The websocket Connection.
	Conn *websocket.Conn `json:"-"`
	// Buffered channel of outbound messages.
	Send chan []byte `json:"-"`
	// ChatRoom that client is registered with
	Room *ChatRoom `json:"-"`
}

// ReadPump pumps messages from the websocket connection to the broker.
//
// The application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a Connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.Room.Broker.CloseClient <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mt, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) || err == io.EOF {
				res, _ := json.Marshal(&ChatEvent{User: c.Username, Msg: fmt.Sprintf("%s has left the room.", c.Username), Color: c.Color})
				c.Room.Broker.Notification <- res
			}
			log.Printf("error: %v", err)
			break
		}
		switch mt {
		case websocket.TextMessage:
			ce, err := validateEvent(data)
			if err != nil {
				log.Printf("Error parsing JSON ChatEvent: %v", err)
				break
			}
			// Set timestamp and room ID
			ce.Timestamp = time.Now()
			ce.RoomID = c.Room.ID

			// Perform requested action
			switch ce.EventType {
			case Unsubscribe:
				// Populate activity
				c.Room.Clients[ce.User].LastActivity = ce.Timestamp
				c.unsubscribe(&ce)
			case Subscribe:
				// LastActivity will be populated in subscribe
				c.subscribe(&ce)
			default:
				// Populate activity
				c.Room.Clients[ce.User].LastActivity = ce.Timestamp
				c.broadcast(&ce)
			}

		default:
			log.Printf("Warning: unknown message type")
		}
	}
}

// WritePump pumps messages from the hub to the websocket Connection.
//
// A goroutine running WritePump is started for each Connection. The
// application ensures that there is at most one writer to a Connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func formatEventData(c *ChatEvent) []byte {
	data, _ := json.Marshal(c)
	return data
}

func (c *Client) broadcast(evt *ChatEvent) {
	evt.EventType = Broadcast
	c.Room.Broker.Notification <- formatEventData(evt)
}

func (c *Client) subscribe(evt *ChatEvent) {
	// Init client values
	c.Username = evt.User
	c.Color = evt.Color
	c.LastActivity = time.Now()
	if err := c.Room.AddClient(c); err != nil {
		log.Println("error adding client:", err.Error())
		return
	}
	log.Println("Adding client to Chatroom: ", evt.User)
	evt.EventType = Subscribe
	evt.Msg = fmt.Sprintf("%s entered the room.", evt.User)
	go func() {
		time.Sleep(200 * time.Millisecond)
		c.Room.Broker.Notification <- formatEventData(evt)
	}()
}

func (c *Client) unsubscribe(evt *ChatEvent) {
	// Remove Client from tracked list
	if err := c.Room.RemoveClient(evt.User); err != nil {
		log.Println("Error removing client", err.Error())
	}
	log.Println(fmt.Sprintf("Unsubscribing %s in room %d", evt.User, c.Room.ID))
	evt.EventType = Unsubscribe
	evt.Msg = fmt.Sprintf("%s left the room.", evt.User)
	go func() {
		time.Sleep(200 * time.Millisecond)
		c.Room.Broker.Notification <- formatEventData(evt)
	}()
}
