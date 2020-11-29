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

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	mux = http.NewServeMux()
	mux.Handle("/chat/", errHandler(handleRoom))
	// If all handlers are desired:
	//registerHandlers()
	//mux = Mux
}

func tearDown() {
}

func TestHandleGetDefaultRoom(t *testing.T) {
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
	cases := []struct {
		title           string
		description     string
		classification  string
		password        string
		expectedOutcome bool
	}{
		{"public room", "this is a public room", "public", "", true},
		{"private room", "this is a private room", "private", "123", true},
		{"secret room", "this is a secret room", "hidden", "!!123abc", true},
		{"Public Chat", "this is a duplicate and should fail", "public", "", false},
	}
	var res data.Outcome
	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			// Refresh writer
			writer = httptest.NewRecorder()
			writer.Header().Set("Content-Type", "application/json")
			// JSON body
			requestJSON := fmt.Sprintf(`{"title":"%s","description":"%s", "classification":"%s", "password":"%s"}`, tc.title, tc.description, tc.classification, tc.password)
			requestBody := strings.NewReader(requestJSON)
			// URI and HTTP method
			request, _ := http.NewRequest("POST", "/chat/", requestBody)
			// Send request
			mux.ServeHTTP(writer, request)
			// Check assertions
			if writer.Code != 200 {
				t.Errorf("Response code is %v", writer.Code)
			}

			json.Unmarshal(writer.Body.Bytes(), &res)
			if res.Success != tc.expectedOutcome {
				t.Fatalf("Unexpected outcome POST chat room \"%s\". Reason: %s", tc.title, res.Error.Msg)
			}
		})
	}
}

func TestHandleGetRooms(t *testing.T) {
	cases := []struct {
		titleOrID           string
		expectedDescription string
		expectedOutcome     bool
	}{
		{"public room", "this is a public room", true},
		{"private room", "this is a private room", true},
		{"secret room", "this is a secret room", true},
		{"0", "This is the default chat, available to everyone!", true},
	}
	var cr data.ChatRoom
	for _, tc := range cases {
		t.Run(tc.titleOrID, func(t *testing.T) {
			// Refresh writer
			writer = httptest.NewRecorder()
			// Craft HTTP req
			writer.Header().Set("Content-Type", "application/json")
			request, _ := http.NewRequest("GET", tc.titleOrID, nil)
			mux.ServeHTTP(writer, request)
			// Check assertions
			if writer.Code != 200 {
				t.Errorf("Response code is %v", writer.Code)
			}

			json.Unmarshal(writer.Body.Bytes(), &cr)
			if strings.ToLower(cr.Title) != tc.titleOrID {
				t.Errorf("Cannot retrieve chat room")
				t.Logf("Response: %s", writer.Body.String())
			} else {
				t.Logf("GET chat room \"%s\"", cr.Title)
			}
		})
	}
}
