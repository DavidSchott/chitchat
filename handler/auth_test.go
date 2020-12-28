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
		roomID                 string
		password               string
		expectedOutcome        bool
		expectedHTTPStatusCode int
	}{
		{PublicTestChatTitle, "", true, 200},
		{PrivateTestChatTitle, "incorrect_pwd", false, 401},
		{PrivateTestChatTitle, "123abc123abc", true, 201},
		{HiddenTestChatTitle, "123abc123abc", true, 201},
		{"does not exist", "123abc123abc", false, 404},
	}

	var result map[string]interface{}
	for _, tc := range cases {
		result = nil
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
				if len(result["token"].(string)) < 100 {
					t.Fatal("Unexpected error generating token", result["token"].(string))
				}
			} else if !tc.expectedOutcome && result["token"] != nil {
				t.Fatal("SECURITY ISSUE: TOKEN UNEXPECTEDLY SET", result)
			}
		})
	}
}

func TestRenewToken(t *testing.T) {
	cases := []struct {
		roomID                 string
		expectedOutcome        bool
		expectedHTTPStatusCode int
	}{
		{PublicTestChatTitle, true, 200},
		{PrivateTestChatTitle, false, 403},
		{PrivateTestChatTitle, true, 201},
		{HiddenTestChatTitle, true, 201},
		{"does not exist", false, 404},
	}
	var result map[string]interface{}
	for _, tc := range cases {
		result = nil
		t.Run(tc.roomID, func(t *testing.T) {
			cr, _ := data.CS.Retrieve(tc.roomID)
			// Refresh writer
			writer = httptest.NewRecorder()
			// URI and HTTP method
			request, _ := http.NewRequest("GET", fmt.Sprintf("/chats/%s/token/renew", tc.roomID), nil)
			request.Header.Set("Content-Type", "application/json")
			if tc.roomID != "does not exist" && cr.Type != data.PublicRoom {
				setJWTHeaders(t, request, tc.roomID, tc.expectedOutcome)
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
			if tc.expectedOutcome && cr.Type != data.PublicRoom {
				if len(result["token"].(string)) < 100 {
					t.Fatal("Unexpected error generating token", result["token"].(string))
				}
			} else if !tc.expectedOutcome && result["token"] != nil {
				t.Fatal("SECURITY ISSUE: TOKEN UNEXPECTEDLY SET", result)
			}
		})
	}
}

// Generates a token and sets it in the request Authorization HTTP header under Bearer scheme
// If intendedValidity is set to false, this will set a faulty token
// This should only be used as a band-aid to keep tests simple and independent for now
func setJWTHeaders(t *testing.T, r *http.Request, id string, intendedValidity bool) {
	t.Helper()
	cr, _ := data.CS.Retrieve(id)
	var myCr *data.ChatRoom = &data.ChatRoom{Password: cr.Password, ID: cr.ID, Title: cr.Title}
	if !intendedValidity {
		myCr.Password = "bogus_incorrect_password"
	}
	tkn, _ := data.EncodeJWT(&data.ChatEvent{User: "test_user", RoomID: cr.ID}, cr, generateUniqueKey(myCr))
	r.Header.Set("Authorization", "Bearer "+tkn)
}
