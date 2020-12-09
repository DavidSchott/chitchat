package handler

import (
	"net/http"
	"strings"

	"github.com/DavidSchott/chitchat/data"
	"github.com/gorilla/mux"
)

// GET /chats
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

// GET /
// Default page
func index(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "layout", "sidebar", "public.header", "index")
}

// GET /chat/<id>/chatbox
// Default page
func chatbox(w http.ResponseWriter, r *http.Request) {
	queries := mux.Vars(r)
	if titleOrID, ok := queries["titleOrID"]; ok {
		cr, err := data.CS.Retrieve(titleOrID)
		if err != nil {
			info("erroneous chats API request", r, err)
		}
		generateHTMLContent(w, &cr, "chat")
	}
}

// GET /chats/<id>/entrance
// Default page
func joinRoom(w http.ResponseWriter, r *http.Request) (err error) {
	queries := mux.Vars(r)
	if titleOrID, ok := queries["titleOrID"]; ok {
		cr, err := data.CS.Retrieve(titleOrID)
		if err != nil {
			info("erroneous chats API request", r, err)
			return err
		}
		generateHTML(w, (strings.ToLower(cr.Type) == data.PrivateRoom || cr.Type == data.HiddenRoom), "layout", "sidebar", "public.header", "entrance")
	}
	return
}
