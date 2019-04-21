package data

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var Clients = make(map[*websocket.Conn]bool)
var Transmission = make(chan *ChatEvent)
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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
}

func (ce ChatEvent) ColorHTML() (css template.CSS) {
	switch ce.Color {
	case "purple":
		css = template.CSS(`background: #7386D5; color: #ffffff !important;`)
	case "blue":
		css = template.CSS(`background: #42a5f5; color: #ffffff !important;`)
	case "red":
		css = template.CSS(`background: #DC143C; color: #ffffff !important;`)
	case "green":
		css = template.CSS(`background: #2E8B57; color: #ffffff !important;`)
	case "gray":
		css = template.CSS(`background: #f1f1f1; color: #000 !important;`)
	case "turquoise":
		css = template.CSS(`background: #40E0D0; color: #000 !important;`)
	case "indigo":
		css = template.CSS(`background: #4B0082; color: #ffffff !important;`)
	case "magenta":
		css = template.CSS(`background: #8B008B; color: #ffffff !important;`)
	case "black":
		css = template.CSS(`background: #000000; color: #ffffff !important;`)
	case "yellow":
		css = template.CSS(`background: #FFD700; color: #000 !important;`)
	case "orange":
		css = template.CSS(`background: #FF8C00; color: #000 !important;`)
	default:
		css = template.CSS(`background: #42a5f5; color: #ffffff !important;`)
	}
	return css
}
