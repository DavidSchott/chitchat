package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/DavidSchott/chitchat/data"
)

const WSHandshakeTimeOut = 45 * time.Second

// TODO: Need to rewrite these tests and/or the design since it cannot be tested without errors (implementation uses unbuffered channels, do not support multiple concurrent read/writes)
func TestHandleWebSocket(t *testing.T) {
	tcs := []struct {
		name            string
		user            string
		eventIterations int
		titleOrID       string
	}{
		{
			name:            "Public Test Chat",
			user:            "Test User",
			eventIterations: 10,
			titleOrID:       "3",
		},
	}

	for _, tt := range tcs {
		// Seed to create random chat messages
		t.Run(tt.name, func(t *testing.T) {
			s, ws := newWSServer(t, tt.titleOrID, router)
			defer s.Close()
			defer ws.Close()
			// Join user to chat room
			joinEvt := data.ChatEvent{EventType: data.Subscribe, User: tt.user}
			expectedEventResponse := joinEvt
			expectedEventResponse.Msg = fmt.Sprintf("%s entered the room.", tt.user)
			compareExpectedActualEvents(t, ws, joinEvt, expectedEventResponse)
			// Send Chat Messages. Need to rewrite this test since design does not support multiple concurrent read/writes
			/*for i := 0; i < tt.eventIterations; i++ {
				msg := fmt.Sprintf("Test message %d for %s from %s", i, tt.name, tt.user)
				sendEvt := data.ChatEvent{EventType: data.Broadcast, User: tt.user, Msg: msg, Color: "Red"}
				compareExpectedActualEvents(t, ws, sendEvt, sendEvt)
			}
			*/
		})
	}
}

func compareExpectedActualEvents(t *testing.T, ws *websocket.Conn, outEvt data.ChatEvent, expectedEvent data.ChatEvent) {
	t.Helper()
	sendWSMessage(t, ws, outEvt)
	t.Log("Sending chat message: " + outEvt.Msg)
	actual := receiveWSMessage(t, ws)
	if actual.Msg != expectedEvent.Msg {
		t.Fatalf("Expected '%+v', got '%+v'", expectedEvent, actual)
	}
}

func newWSServer(t *testing.T, titleOrID string, h http.Handler) (*httptest.Server, *websocket.Conn) {
	t.Helper()
	s := httptest.NewServer(h)
	// Transform URL from HTTP to wss://
	wsURL := httpToWS(t, s.URL)
	wsURL = wsURL + fmt.Sprintf("/chats/%s/ws", titleOrID)
	d := websocket.Dialer{ReadBufferSize: upgrader.ReadBufferSize, WriteBufferSize: upgrader.WriteBufferSize, HandshakeTimeout: WSHandshakeTimeOut, Proxy: http.ProxyFromEnvironment}
	// Open WebSocket Conn
	ws, resp, err := d.Dial(wsURL, nil)
	//ws, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Check we have upgraded from HTTP to WebSocket protocol successfully
	if got, want := resp.StatusCode, http.StatusSwitchingProtocols; got != want {
		t.Errorf("resp.StatusCode = %q, want %q", got, want)
	}

	return s, ws
}

func sendWSMessage(t *testing.T, ws *websocket.Conn, ce data.ChatEvent) {
	t.Helper()

	m, err := json.Marshal(ce)
	if err != nil {
		t.Fatal(err)
	}
	if err := ws.WriteMessage(websocket.TextMessage, m); err != nil {
		t.Fatalf("%v", err)
	}
}

func receiveWSMessage(t *testing.T, ws *websocket.Conn) data.ChatEvent {
	t.Helper()

	_, m, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	var reply data.ChatEvent
	reply, err = data.ValidateEvent(m)
	//err = json.Unmarshal(m, &reply)
	if err != nil {
		t.Fatal(err)
	}
	return reply
}

func httpToWS(t *testing.T, u string) string {
	t.Helper()

	wsURL, err := url.Parse(u)
	if err != nil {
		t.Fatal(err)
	}

	switch wsURL.Scheme {
	case "http":
		wsURL.Scheme = "ws"
	case "https":
		wsURL.Scheme = "wss"
	}

	return wsURL.String()
}
