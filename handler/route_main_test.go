package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/DavidSchott/chitchat/data"
)

var writer *httptest.ResponseRecorder
var mux *http.ServeMux
var chatRoomTitle string

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	mux = http.NewServeMux()
	mux.Handle("/chat/", errHandler(handleRoom))
	chatRoomTitle = "test chat room"
	// If all handlers are desired:
	//registerHandlers()
	//mux = Mux
}

func tearDown() {
}

func TestHandleGet(t *testing.T) {
	writer = httptest.NewRecorder()
	writer.Header().Set("Content-Type", "application/json")
	request, _ := http.NewRequest("GET", "/chat/1", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	var cr data.ChatRoom
	json.Unmarshal(writer.Body.Bytes(), &cr)
	if cr.ID != 1 {
		t.Errorf("Cannot retrieve chat room")
		t.Logf("Response: %s", writer.Body.String())
	} else {
		t.Logf("GET chat room \"%s\"", cr.Title)
	}
}

func TestHandlePost(t *testing.T) {
	writer = httptest.NewRecorder()
	writer.Header().Set("Content-Type", "application/json")
	requestJSON := fmt.Sprintf(`{"title":"%s","description":"This is a test", "classification":"public"}`, chatRoomTitle)
	requestBody := strings.NewReader(requestJSON)
	request, _ := http.NewRequest("POST", "/chat/", requestBody)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	var res data.Outcome
	json.Unmarshal(writer.Body.Bytes(), &res)
	if !res.Success {
		t.Errorf("POST chat room \"%s\" FAILED. Reason: %s", chatRoomTitle, res.Error)
	} else {
		t.Logf("POST chat room \"%s\"", chatRoomTitle)
	}

}
