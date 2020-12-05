package handler

import (
	"encoding/json"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/DavidSchott/chitchat/data"
	"github.com/gorilla/mux"
)

const (
	sessionKey string = "secret_cookie"
)

// TODO: Implement as a chained handler
func authorize(h errHandler) errHandler {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		// Skip authorization for special case of GET /chats/<id>
		if name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name(); strings.HasSuffix(name, "handleRoom") && r.Method == http.MethodGet {
			return h(w, r)
		}
		queries := mux.Vars(r)
		if titleOrID, ok := queries["titleOrID"]; ok {
			cr, err := data.CS.Retrieve(titleOrID)
			if err != nil {
				info("erroneous chats API request", r, err)
				return err
			}
			cookieSecret, err := r.Cookie(sessionKey)
			if cr.Type != data.PublicRoom && err != nil || cr.Type != data.PublicRoom && !cr.MatchesPassword(cookieSecret.Value) {
				return &data.APIError{
					Code:  104,
					Field: "password",
				}
			}
			return h(w, r)
		}
		return
	}
}

// POST /chats/{titleOrID}/token
func login(w http.ResponseWriter, r *http.Request) (err error) {
	w.Header().Set("Content-Type", "application/json")
	// read in request
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var c data.ChatEvent
	if err := json.Unmarshal(body, &c); err != nil {
		danger("Error parsing token request", r)
	}
	queries := mux.Vars(r)
	if titleOrID, ok := queries["titleOrID"]; ok {
		cr, err := data.CS.Retrieve(titleOrID)
		if err != nil {
			info("erroneous chats API request", r, err)
			return err
		}
		if cr.Type == data.PublicRoom {
			// Ignore public room
			ReportStatus(w, true, nil)
		} else if cr.MatchesPassword(c.Password) {
			// Success! Set Password
			cookieSecret := http.Cookie{
				Name:     sessionKey,
				Value:    c.Password,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			}
			http.SetCookie(w, &cookieSecret)
			ReportStatus(w, true, nil)
		} else {
			return &data.APIError{
				Code:  304,
				Field: "secret",
			}
		}
	}
	return
}
