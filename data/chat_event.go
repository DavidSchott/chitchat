package data

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	// Subscribe is used to broadcast a message indicating user has joined ChatRoom
	Subscribe = "join"
	// Broadcast is used to broadcast messages to all subscribed users
	Broadcast = "send"
	// Unsubscribe is used to broadcast a message indicating user has left ChatRoom
	Unsubscribe = "leave"
)

// ChatEvent represents a message event in an associated ChatRoom
type ChatEvent struct {
	EventType string    `json:"event_type,omitempty"`
	User      string    `json:"name,omitempty"`
	RoomID    string    `json:"room_id,omitempty"`
	Color     string    `json:"color,omitempty"`
	Msg       string    `json:"msg,omitempty"`
	Password  string    `json:"secret,omitempty"`
	Timestamp time.Time `json:"time,omitempty"`
}

// ValidateEvent ensures data is a valid JSON representation of Chat Event and can be parsed as such
func ValidateEvent(data []byte) (ChatEvent, error) {
	var evt ChatEvent

	if err := json.Unmarshal(data, &evt); err != nil {
		return evt, &APIError{Code: 303}
	}

	if evt.User == "" {
		return evt, &APIError{Code: 303, Field: "name"}
	} else if evt.Msg == "" && strings.ToLower(evt.EventType) == Broadcast {
		return evt, &APIError{Code: 303, Field: "msg"}
	}

	return evt, nil
}
