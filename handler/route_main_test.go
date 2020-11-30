package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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

func TestHandlePost(t *testing.T) {
	cases := []struct {
		title                  string
		description            string
		classification         string
		password               string
		expectedOutcome        bool
		expectedHTTPStatusCode int
	}{
		{"public room", "this is a public room", "public", "", true, 201},
		{"private room", "this is a private room", "private", "123", true, 201},
		{"secret room", "this is a secret room", "hidden", "!!123abc", true, 201},
		{"Public Chat", "this is a duplicate and should fail", "public", "", false, 400},
	}
	var res data.ChatRoom
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
			if writer.Code != tc.expectedHTTPStatusCode {
				t.Errorf("Response code is %v", writer.Code)
			}

			json.Unmarshal(writer.Body.Bytes(), &res)
			assertion := assertTrue(res.Title == tc.title, res.Description == tc.description, res.Type == tc.classification)
			// TODO: Check all fields
			if assertion != tc.expectedOutcome {
				t.Fatal("Unexpected result POST chat room. Response: ", res)
			}
		})
	}
}

func TestHandleGetRooms(t *testing.T) {
	cases := []struct {
		titleOrID              string
		expectedDescription    string
		expectedHTTPStatusCode int
		expectedOutcome        bool
	}{
		{"1", "This is the default chat, available to everyone!", 200, true},
		{"public room", "this is a public room", 200, true},
		{"private room", "this is a private room", 200, true},
		{"secret room", "this is a secret room", 200, true},
		{"this room does not exist", "this is a problem", 404, false},
		{"", "this is a bad request", 404, false},
	}
	var cr data.ChatRoom
	var failOutcome data.Outcome
	var matchConditions bool
	for _, tc := range cases {
		t.Run(tc.titleOrID, func(t *testing.T) {
			// Refresh writer
			writer = httptest.NewRecorder()
			// Craft HTTP req
			writer.Header().Set("Content-Type", "application/json")
			request, _ := http.NewRequest("GET", fmt.Sprintf("/chat/%s", tc.titleOrID), nil)
			mux.ServeHTTP(writer, request)
			// Check assertions
			if writer.Code != tc.expectedHTTPStatusCode {
				t.Errorf("Response code is %v", writer.Code)
			}
			if tc.expectedOutcome {
				json.Unmarshal(writer.Body.Bytes(), &cr)
				matchConditions = strings.ToLower(cr.Title) == tc.titleOrID || strconv.Itoa(cr.ID) == tc.titleOrID
			} else {
				err := json.Unmarshal(writer.Body.Bytes(), &failOutcome)
				if err != nil {
					t.Log(err.Error())
					matchConditions = false
				}
				// Check error return code is as expected
				matchConditions = assertTrue(failOutcome.Error.Code == 101, !failOutcome.Success)
			}
			// If assumed test checks fail
			if !matchConditions {
				t.Errorf("Unexpected result during GET chat room %s", tc.titleOrID)
				t.Logf("Response: %s", writer.Body.String())
			}
		})
	}
}

func assertTrue(vals ...bool) bool {
	allTrue := true
	for _, val := range vals {
		allTrue = allTrue && val
	}
	return allTrue
}

/*
func TestHandlePut(t *testing.T) {
	cases := []struct {
		titleOrID       string
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
		t.Run(tc.titleOrID, func(t *testing.T) {
			// Refresh writer
			writer = httptest.NewRecorder()
			writer.Header().Set("Content-Type", "application/json")
			// JSON body
			requestJSON := fmt.Sprintf(`{"title":"%s","description":"%s", "classification":"%s", "password":"%s"}`, tc.titleOrID, tc.description, tc.classification, tc.password)
			requestBody := strings.NewReader(requestJSON)
			// URI and HTTP method
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/chat/%s", tc.titleOrID), requestBody)
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
*/
