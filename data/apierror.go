package data

import (
	"fmt"
	"strconv"
)

/* APIError returns a JSON-string with a code and error message for chat-related API calls
 0 = success
-10* = room errors
  -101 = room not found
  -102 = room duplicate
  -103 = invalid json spec
  -104 = unauthorized
  -105 = invalid content (e.g. short password)

-20* = client errors
  -201 = client not found
  -202 = client duplicate
  -203 = invalid json spec
  -204 = unauthorized

-30* = event errors
  - 301 = could not establish session
  - 303 = invalid json spec
  - 304 = unauthorized
  - 305 = unsupported client device
  - 306 = invalid signing method
-40* = token errors
  - 401 = Invalid signature
  - 402 = Unauthorized signing method
  - 403 = Invalid token
*/

// APIError represents an error that was thrown at some point with some relevant information for users to correct their input
type APIError struct {
	Code  int    `json:"code,omitempty"`
	Msg   string `json:"error,omitempty"`
	Field string `json:"field,omitempty"`
}

// Outcome represents status indicating success/failure of an operation
type Outcome struct {
	Status bool      `json:"status"`
	Error  *APIError `json:"error,omitempty"`
}

// SetMsg will set the Msg based on the provided code
func (e *APIError) SetMsg() {
	switch e.Code {
	case 101:
		e.Msg = "Room error: Room not found"
	case 102:
		e.Msg = "Room error: Duplicate room"
	case 103:
		e.Msg = "Room error: Invalid JSON"
	case 104:
		e.Msg = "Room error: Unauthorized operation"
	case 105:
		e.Msg = "Room error: Invalid content"
	case 201:
		e.Msg = "Client error: User not found"
	case 202:
		e.Msg = "Client error: Duplicate username"
	case 203:
		e.Msg = "Client error: Invalid JSON"
	case 204:
		e.Msg = "Client error: Unauthorized operation"
	case 301:
		e.Msg = "Could not establish session"
	case 303:
		e.Msg = "Invalid JSON"
	case 304:
		e.Msg = "Unauthorized operation"
	case 305:
		e.Msg = "Unsupported client device"
	case 401:
		e.Msg = "Token error: Invalid signature"
	case 402:
		e.Msg = "Token error: Unauthorized signing method"
	case 403:
		e.Msg = "Token error: Invalid token"
	default:
		e.Msg = "Unknown error: " + e.Msg
	}
}

func (e *APIError) Error() string {
	e.SetMsg()
	if e.Field != "" {
		return fmt.Sprintf("{\"error\": \"%s\", \"code\": %d, \"field\": \"%s\"}", e.Msg, e.Code, e.Field)
	}
	return fmt.Sprintf("{\"error\": \"%s\", \"code\": %d}", e.Msg, e.Code)
}

func isInt(titleorID string) int {
	if id, err := strconv.Atoi(titleorID); err == nil {
		return id
	}
	return -1
}
