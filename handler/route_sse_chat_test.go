package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DavidSchott/chitchat/data"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		RoomID                 string
		password               string
		expectedOutcome        bool
		expectedHTTPStatusCode int
		expectedAPIErrorCode   int
	}{
		{"1", "", true, 200, 0},
	}
	var res data.Outcome
	for _, tc := range cases {
		t.Run(tc.RoomID, func(t *testing.T) {
			// Refresh writer
			writer = httptest.NewRecorder()
			// JSON body
			requestJSON := fmt.Sprintf(`{"room_id":"%s","secret":"%s"}`, tc.RoomID, tc.password)
			requestBody := strings.NewReader(requestJSON)
			// URI and HTTP method
			request, _ := http.NewRequest("POST", fmt.Sprintf("/chats/%s/token", tc.RoomID), requestBody)
			request.Header.Set("Content-Type", "application/json")
			// Send request
			router.ServeHTTP(writer, request)
			// Check assertions
			if writer.Code != tc.expectedHTTPStatusCode {
				t.Errorf("Response code is %v", writer.Code)
			}
			if err := json.Unmarshal(writer.Body.Bytes(), &res); err != nil {
				t.Fatal("Unexpected result authorizing. Response: ", writer.Body.String())
			}
			if res.Status != tc.expectedOutcome {
				t.Error("Unexpected result authorizing. Response: ", writer.Body.String())
			}
		})
	}
}
