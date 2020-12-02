package handler

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"

	"github.com/DavidSchott/chitchat/data"
)

const (
	sessionKey string = "secret_cookie"
)

// GET /
// Default page
func index(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "layout", "sidebar", "public.header", "index")
}

// GET /about
// Default page
func about(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "layout", "sidebar", "public.header", "about")
}

// GET /test
// Default page
func test(w http.ResponseWriter, r *http.Request) (err error) {
	generateHTML(w, "", "layout", "sidebar", "public.header", "test")
	return
}

// GET /chat/join/<id>
// Default page
func joinRoom(w http.ResponseWriter, r *http.Request) (err error) {
	//ID, err := strconv.Atoi(path.Base(r.URL.Path))
	var ID string = path.Base(r.URL.Path)
	info("joining room", ID)
	cr, err := data.CS.Retrieve(ID)
	if err != nil {
		return err
	}
	generateHTML(w, (strings.ToLower(cr.Type) == data.PrivateRoom || cr.Type == data.HiddenRoom), "layout", "sidebar", "public.header", "entrance")
	return
}

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

// GET /chat/list
func listChats(w http.ResponseWriter, r *http.Request) {
	rooms, err := data.CS.Chats()
	if err != nil {
		errorMessage(w, r, "Cannot retrieve chats")
	} else {
		// to return back to refreshing page:
		//generateHTML(w, &rooms, "layout", "sidebar", "public.header", "list")
		generateHTMLContent(w, &rooms, "list")
		return
	}
}

// GET /chat/box/<id>
// Default page
func chatbox(w http.ResponseWriter, r *http.Request) {
	//ID, err := strconv.Atoi(path.Base(r.URL.Path))
	var ID string = path.Base(r.URL.Path)
	info("joining room", ID)
	cr, err := data.CS.Retrieve(ID)
	if err != nil {
		//w.Write([]byte(err.Error()))
		//return
		p(err.Error())
	} else {
		generateHTMLContent(w, &cr, "chat")
		return
	}
}

// main handler function
func handleRoom(w http.ResponseWriter, r *http.Request) (err error) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		err = handleGet(w, r)
	case "POST":
		err = handlePost(w, r)
	case "PUT":
		err = handlePut(w, r)
	case "DELETE":
		err = handleDelete(w, r)
	}
	return err
}

// Retrieve a chat room
// GET /chat/1
func handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	title := path.Base(r.URL.Path)
	cr, err := data.CS.Retrieve(title)
	if err != nil {
		info("error getting chat room: " + title)
		return err
	}
	res, err := cr.ToJSON()
	info("retrieved chat room:", cr.Title)
	w.Write(res)
	return
}

// Create a ChatRoom
// POST /chat/
func handlePost(w http.ResponseWriter, r *http.Request) (err error) {
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
func handlePut(w http.ResponseWriter, r *http.Request) (err error) {
	title := path.Base(r.URL.Path)
	var cr data.ChatRoom
	currentChatRoom, err := data.CS.Retrieve(title)
	if err != nil {
		return
	}
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
			warning("error attempting to authorize "+title+" by PUT:", cr)
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
func handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	title := path.Base(r.URL.Path)
	cr, err := data.CS.Retrieve(title)
	if err != nil {
		return
	}
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
