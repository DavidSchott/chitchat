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
	RoomID    int       `json:"id"`
	Color     string    `json:"color"`
	Msg       string    `json:"msg"`
}
