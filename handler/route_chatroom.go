package handler

import (
	"encoding/json"
	"net/http"

	"github.com/DavidSchott/chitchat/data"
	"github.com/gorilla/mux"
)

const (
	sessionKey string = "secret_cookie"
)

// NOT USED. GET /chat/join/<id>
// TODO: Implement as a chained handler
func authorize(h errHandler) errHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		cookieSecret, err := r.Cookie("secret_cookie")
		if err != nil {
			return &data.APIError{
				Code:  304,
				Field: "secret",
			}
		}
		// TODO: Actually check this is valid
		if cookieSecret.Value != "password" {
			return &data.APIError{
				Code:  304,
				Field: "secret",
			}
		}
		return h(w, r)
	}

}

// main handler function
func handleRoom(w http.ResponseWriter, r *http.Request) (err error) {
	queries := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	if titleOrID, ok := queries["titleOrID"]; ok {
		cr, err := data.CS.Retrieve(titleOrID)
		if err != nil {
			info("erroneous chats API request", r, err)
			return err
		}
		switch r.Method {
		case "GET":
			err = handleGet(w, r, cr)
			return err
		case "PUT":
			err = handlePut(w, r, cr, titleOrID)
			return err
		case "DELETE":
			err = handleDelete(w, r, cr)
			return err
		}
	} else {
		err = &data.APIError{
			Code: 103,
		}
	}

	return err
}

// Retrieve a chat room
// GET /chat/1
func handleGet(w http.ResponseWriter, r *http.Request, cr *data.ChatRoom) (err error) {
	res, err := cr.ToJSON()
	if err != nil {
		return
	}
	info("retrieved chat room:", cr.Title)
	w.Write(res)
	return
}

// Create a ChatRoom
// POST /chats
func handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	w.Header().Set("Content-Type", "application/json")
	// read in request
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	// create ChatRoom obj
	var cr data.ChatRoom
	if err = json.Unmarshal(body, &cr); err != nil {
		warning("error encountered reading POST:", err.Error())
		return err
	}
	if err = data.CS.Add(&cr); err != nil {
		warning("error encountered adding chat room:", err.Error())
		return err
	}
	// Retrieve updated object
	createdChatRoom, err := data.CS.Retrieve(cr.Title)
	if err != nil {
		return err
	}
	response, _ := createdChatRoom.ToJSON()
	w.WriteHeader(201)
	w.Write(response)
	return
}

// Update a room
// PUT /chat/<id>
func handlePut(w http.ResponseWriter, r *http.Request, currentChatRoom *data.ChatRoom, title string) (err error) {
	var cr data.ChatRoom
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	if err = json.Unmarshal(body, &cr); err != nil {
		warning("error encountered updating chat room:", err.Error())
		return
	}
	// Authorize
	if currentChatRoom.Type != data.PublicRoom {
		// if isn't public room, need to authorize
		cookieSecret, err := r.Cookie("secret_cookie")
		if err != nil {
			warning("error attempting to authorize ", r, " PUT:", r)
			return &data.APIError{
				Code:  104,
				Field: "password",
			}
		}
		if cookieSecret.Value != currentChatRoom.Password {
			return &data.APIError{
				Code:  104,
				Field: "password",
			}
		}
	}
	if err = data.CS.Update(title, &cr); err != nil {
		warning("error encountered updating chat room:", cr, err.Error())
		return
	}
	// Retrieve updated object
	modifiedChatRoom, err := data.CS.RetrieveID(currentChatRoom.ID)
	if err != nil {
		return err
	}
	info("updated chat room:", title)
	response, _ := modifiedChatRoom.ToJSON()
	w.Write(response)
	return
}

// Delete a room
// DELETE /chat/<id>
func handleDelete(w http.ResponseWriter, r *http.Request, cr *data.ChatRoom) (err error) {
	// TODO: authorize
	err = data.CS.Delete(cr)
	if err != nil {
		warning("error encountered deleting chat room:", err.Error())
		ReportStatus(w, false, err.(*data.APIError))
		return
	}
	// report on status
	info("deleted chat room:", cr.Title)
	ReportStatus(w, true, nil)
	return
}
