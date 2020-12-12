package data

import (
	"encoding/json"
	"time"
)

const (
	// Subscribe is used to broadcast a message indicating user has join room
	Subscribe = "join"
	// Broadcast is used to broadcast messages to all subscribed users
	Broadcast = "send"
	// Unsubscribe is used to broadcast a message indicating user has left room
	Unsubscribe = "leave"
)

// ChatEvent represents a message event in an associated chat room
type ChatEvent struct {
	EventType string    `json:"event_type,omitempty"`
	User      string    `json:"name,omitempty"`
	RoomID    int       `json:"room_id,omitempty"`
	Color     string    `json:"color,omitempty"`
	Msg       string    `json:"msg,omitempty"`
	Password  string    `json:"secret,omitempty"`
	Timestamp time.Time `json:"time,omitempty"`
}

// validateEvent so that we know it's a valid JSON representation of Chat event
func validateEvent(data []byte) (ChatEvent, error) {
	var evt ChatEvent

	if err := json.Unmarshal(data, &evt); err != nil {
		return evt, &APIError{Code: 303}
	}

	if evt.User == "" || evt.Msg == "" {
		return evt, &APIError{Code: 303, Field: "name|msg"}
	}

	return evt, nil
}
