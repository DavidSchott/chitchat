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
	EventType string    `json:"event_type,omitempty"`
	User      string    `json:"name,omitempty"`
	RoomID    int       `json:"room_id,omitempty"`
	Color     string    `json:"color,omitempty"`
	Msg       string    `json:"msg,omitempty"`
	Password  string    `json:"secret,omitempty"`
	Timestamp time.Time `json:"time,omitempty"`
}