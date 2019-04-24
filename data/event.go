package data

import (
	"time"
)

const (
	Subscribe   = "join"
	Broadcast   = "send"
	Unsubscribe = "leave"
)

type ChatEvent struct {
	EventType string    `json:"type"`
	User      string    `json:"name"`
	Timestamp time.Time `json:"time"`
	RoomID    int       `json:"id,omitempty"`
	Color     string    `json:"color,omitempty"`
	Msg       string    `json:"msg,omitempty"`
}
