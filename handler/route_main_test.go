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
		visibility             string
		password               string
		expectedOutcome        bool
		expectedHTTPStatusCode int
		expectedAPIErrorCode   int
	}{
		{"public room", "this is a public room", "public", "", true, 201, 0},
		{"private room", "this is a private room", "private", "password123", true, 201, 0},
		{"secret room", "this is a secret room", "hidden", "!!123abcpassword", true, 201, 0},
		{"Public Chat", "this is a duplicate of default room and should fail", "public", "", false, 400, 102},
		{"", "the title shall not be empty", "public", "", false, 400, 105},
		{"bad room", "the visibility shall not be empty", "", "", false, 400, 105},
		{"bad private room", "password shall not be too short for private rooms", "private", "", false, 400, 105},
		{"bad hidden room", "password shall not be too short  for hidden rooms", "hidden", "", false, 400, 105},
		{"weird public room", "passwords given to a public room shall fail to avoid accidents", "public", "badpwd", false, 400, 105},
	}
	var res data.ChatRoom
	var failedOutcome data.Outcome
	var matchConditions bool
	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			// Refresh writer
			writer = httptest.NewRecorder()
			// JSON body
			requestJSON := fmt.Sprintf(`{"title":"%s","description":"%s", "visibility":"%s", "password":"%s"}`, tc.title, tc.description, tc.visibility, tc.password)
			requestBody := strings.NewReader(requestJSON)
			// URI and HTTP method
			request, _ := http.NewRequest("POST", "/chat/", requestBody)
			request.Header.Set("Content-Type", "application/json")
			// Send request
			mux.ServeHTTP(writer, request)
			// Check assertions
			if writer.Code != tc.expectedHTTPStatusCode {
				t.Errorf("Response code is %v", writer.Code)
			}
			if tc.expectedOutcome {
				json.Unmarshal(writer.Body.Bytes(), &res)
				matchConditions = assertTrue(res.Title == tc.title, res.Description == tc.description, res.Type == tc.visibility, res.ID > 1)
			} else {
				json.Unmarshal(writer.Body.Bytes(), &failedOutcome)
				matchConditions = assertTrue(!failedOutcome.Success, failedOutcome.Error.Code == tc.expectedAPIErrorCode)
			}

			// TODO: Check all fields
			if !matchConditions {
				t.Fatal("Unexpected result POST chat room. Response: ", writer.Body.String())
			}
		})
	}
}

func TestHandleGet(t *testing.T) {
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
			// Refresh writer TODO: Recycle old one instead.
			writer = httptest.NewRecorder()
			// Craft HTTP req
			request, _ := http.NewRequest("GET", fmt.Sprintf("/chat/%s", tc.titleOrID), nil)
			request.Header.Set("Content-Type", "application/json")
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

func TestHandlePut(t *testing.T) {
	cases := []struct {
		titleOrID              string
		title                  string
		description            string
		visibility             string
		password               string
		expectedOutcome        bool
		expectedHTTPStatusCode int
		expectedAPIErrorCode   int
	}{
		{"1", "default chat room", "renamed", "public", "", true, 200, 0},
		{"public room", "public chat renamed", "renamed", "public", "", true, 200, 0},
		{"3", "private room renamed", "renamed", "private", "password123", true, 200, 0},
		{"3", "private room renamed failure", "bad password", "private", "incorrectpassword", false, 401, 104},
		{"4", "hidden room renamed", "renamed", "hidden", "!!123abcpassword", true, 200, 0},
	}
	var res data.ChatRoom
	var failedOutcome data.Outcome
	var matchConditions bool
	for _, tc := range cases {
		t.Run(tc.titleOrID, func(t *testing.T) {
			// Refresh writer
			writer = httptest.NewRecorder()
			// JSON body
			requestJSON := fmt.Sprintf(`{"title":"%s","description":"%s", "visibility":"%s"}`, tc.title, tc.description, tc.visibility)
			requestBody := strings.NewReader(requestJSON)
			// URI and HTTP method
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/chat/%s", tc.titleOrID), requestBody)
			request.Header.Set("Content-Type", "application/json")
			if len(tc.password) > 1 {
				request.Header.Set("Cookie", fmt.Sprintf("secret_cookie=%s", tc.password))
			}
			// Send request
			mux.ServeHTTP(writer, request)
			// Check assertions
			if writer.Code != tc.expectedHTTPStatusCode {
				t.Errorf("Unexpected response code is %v", writer.Code)
			}
			if tc.expectedOutcome {
				json.Unmarshal(writer.Body.Bytes(), &res)
				matchConditions = assertTrue(res.Title == tc.title, res.Description == tc.description, res.Type == tc.visibility, res.ID >= 1)
			} else {
				json.Unmarshal(writer.Body.Bytes(), &failedOutcome)
				matchConditions = assertTrue(!failedOutcome.Success, failedOutcome.Error.Code == tc.expectedAPIErrorCode)
			}

			// TODO: Check all fields
			if !matchConditions {
				t.Fatal("Unexpected result PUT chat room: ", writer.Body.String())
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
