package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		roomID                 string
		password               string
		expectedOutcome        bool
		expectedHTTPStatusCode int
	}{
		{"1", "", true, 200},
		{"2", "incorrect_pwd", false, 401},
		{"2", "123abc123abc", true, 201},
		{"hidden chat", "123abc123abc", true, 201},
		{"does not exist", "123abc123abc", false, 404},
	}

	var result map[string]interface{}
	for _, tc := range cases {
		t.Run(tc.roomID, func(t *testing.T) {
			// Refresh writer
			writer = httptest.NewRecorder()
			// JSON body
			requestJSON := fmt.Sprintf(`{"room_id":"%s","secret":"%s", "name":"test_user"}`, tc.roomID, tc.password)
			requestBody := strings.NewReader(requestJSON)
			// URI and HTTP method
			request, _ := http.NewRequest("POST", fmt.Sprintf("/chats/%s/token", tc.roomID), requestBody)
			request.Header.Set("Content-Type", "application/json")
			// Send request
			router.ServeHTTP(writer, request)
			// Check assertions
			// HTTP Status Code
			if writer.Code != tc.expectedHTTPStatusCode {
				t.Errorf("Response code is %v", writer.Code)
			}
			// Error parsing
			if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil {
				t.Fatal("Unexpected result authorizing. Response: ", writer.Body.String())
			}
			// Actual outcome is as expected
			if result["status"] != tc.expectedOutcome {
				t.Error("Unexpected result authorizing. Response: ", writer.Body.String())
			}
			// Check that token is set
			if tc.expectedOutcome && tc.password != "" {
				if len(result["token"].(string)) != 152 {
					t.Fatal("Unexpected error generating token", result["token"].(string))
				}
			}
		})
	}
}

func TestRenewToken(t *testing.T) {
	cases := []struct {
		roomID                 string
		password               string
		expectedOutcome        bool
		expectedHTTPStatusCode int
	}{
		{"1", "", true, 200},
		{"2", "incorrect_pwd", false, 403},
		{"2", "123abc123abc", true, 201},
		{"hidden chat", "123abc123abc", true, 201},
		{"does not exist", "123abc123abc", false, 404},
	}
	var result map[string]interface{}
	for _, tc := range cases {
		t.Run(tc.roomID, func(t *testing.T) {
			// Refresh writer
			writer = httptest.NewRecorder()
			// URI and HTTP method
			request, _ := http.NewRequest("GET", fmt.Sprintf("/chats/%s/token/renew", tc.roomID), nil)
			request.Header.Set("Content-Type", "application/json")
			if tc.roomID != "does not exist" {
				setJWTHeaders(request, tc.password, tc.roomID)
			}

			// Send request
			router.ServeHTTP(writer, request)
			// Check assertions
			if writer.Code != tc.expectedHTTPStatusCode {
				t.Errorf("Response code is %v", writer.Code)
			}
			if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil {
				t.Fatal("Unexpected result authorizing. Response: ", writer.Body.String())
			}
			if result["status"] != tc.expectedOutcome {
				t.Error("Unexpected result authorizing. Response: ", writer.Body.String())
			}
			// Check that token is set
			if tc.expectedOutcome && tc.password != "" {
				if len(result["token"].(string)) != 152 {
					t.Fatal("Unexpected error generating token", result["token"].(string))
				}
			}
		})
	}
}
