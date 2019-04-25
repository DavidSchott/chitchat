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
  -104 = invalid content (e.g. short password)

-20* = client errors
  -201 = client not found
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
		e.msg = "Room error: Invalid content"
	case 201:
		e.msg = "Client error: User not found"
	case 202:
		e.msg = "Client error: Duplicate username"
	case 203:
		e.msg = "Client error: Invalid JSON"
	case 204:
		e.msg = "Client error: Unauthorized access"
	case 303:
		e.msg = "Event error: Invalid JSON"
	case 304:
		e.msg = "Event error: Unauthorized access"
	default:
		e.msg = "Unknown error: " + e.msg
	}
	if e.field != "" {
		return fmt.Sprintf("{\"msg\": \"%s\", \"code\": %d, \"field\": \"%s\"}", e.msg, e.code, e.field)
	}
	return fmt.Sprintf("{\"msg\": \"%s\", \"code\": %d}", e.msg, e.code)
}

func isInt(titleorID string) int {
	if id, err := strconv.Atoi(titleorID); err == nil {
		return id
	}
	return -1
}
