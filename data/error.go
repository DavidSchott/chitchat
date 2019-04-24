package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type appHandler func(http.ResponseWriter, *http.Request) error

/*
 0 = success
-10* = room errors
  -101 = room not found
  -102 = room duplicate
  -103 = invalid json spec
  -104 = invalid content (e.g. short password)

-20* = client errors
  -201 = client not found // TODO: remove?
  -202 = client duplicate
  -203 = invalid json spec
  -204 = unauthorized

-30* = event errors
  - 303 = invalid json spec
  - 304 = unauthorized
*/
type APIError struct {
	code  int    `json:"code"`
	msg   string `json:"msg"`
	field string `json:"field"`
}

func (e *APIError) Error() string {
	switch e.code {
	case 101:
		e.msg = "Room error: Room not found"
	case 102:
		e.msg = "Room error: Duplicate room"
	case 103:
		e.msg = "Room error: Invalid JSON"
	case 104:
		e.msg = "Room error: Invalid content for field " + e.field
	case 201:
		e.msg = "Client error: User not found"
	case 202:
		e.msg = "Client error: Duplicate username"
	case 203:
		e.msg = "Client error: Invalid JSON for field " + e.field
	case 204:
		e.msg = "Client error: Unauthorized access"
	case 303:
		e.msg = "Event error: Invalid JSON for field " + e.field
	case 304:
		e.msg = "Event error: Unauthorized access"
	default:
		e.msg = "Unknown error: " + e.msg
	}
	return fmt.Sprintf("{\"msg\": \"%s\", \"code\": %d}", e.msg, e.code)
	//return e.err
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		if apierr, ok := err.(*APIError); ok {
			w.Header().Set("Content-Type", "application/json")
			json, _ := json.Marshal(apierr)
			w.Write(json)
			//warning("API error:", apierr.Error())
		} else {
			//danger("Server error", err.Error())
			http.Error(w, err.Error(), 500)
		}
	}
}

func isInt(titleorID string) int {
	if id, err := strconv.Atoi(titleorID); err == nil {
		return id
	}
	return -1
}
