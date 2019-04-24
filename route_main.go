package main

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	"github.com/DavidSchott/chitchat/data"
)

// GET /
// Default page
func index(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "layout", "sidebar", "public.header", "index")
}

// GET /test
// Default page
func test(w http.ResponseWriter, r *http.Request) {
	ce := data.ChatEvent{
		Color: "black",
	}
	generateHTML(w, &ce, "layout", "sidebar", "public.header", "chat")
}

// GET /chat/join/<id>
// Default page
func joinRoom(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(path.Base(r.URL.Path))
	info("joining room", ID)
	cr, err := data.CS.RetrieveID(ID)
	if err != nil {
		return
	}
	generateHTML(w, &cr, "layout", "sidebar", "public.header", "chat")
}

// GET /chat/list
func listChats(w http.ResponseWriter, r *http.Request) {
	rooms, err := data.CS.Chats()
	if err != nil {
		error_message(w, r, "Cannot retrieve chats")
	} else {
		// to return back to refreshing page:
		//generateHTML(w, &rooms, "layout", "sidebar", "public.header", "list")
		generateHTMLContent(w, &rooms, "list")
		return
	}
}

// main handler function
func handleRoom(w http.ResponseWriter, r *http.Request) {
	var err error
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Retrieve a chat room
// GET /chat/1
func handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	title := path.Base(r.URL.Path)
	cr, err := data.CS.Retrieve(title)
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(&cr, "", "\t\t")
	if err != nil {
		return
	}
	// report on success
	info("retrieved chat room:", cr.Title)
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
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
	json.Unmarshal(body, &cr)
	err = data.CS.Add(&cr)
	// report on success/error
	if err != nil {
		warning("error encountered creating chat room:", err.Error())
		ReportSuccess(w, false, err.Error())
		return
	}
	info("created chat room:", cr.Title)
	ReportSuccess(w, true, "")
	return
}

// Update a room
// PUT /chat/<id>
func handlePut(w http.ResponseWriter, r *http.Request) (err error) {
	title := path.Base(r.URL.Path)
	cr, err := data.CS.Retrieve(title)
	if err != nil {
		return
	}
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	json.Unmarshal(body, &cr)
	err = data.CS.Update(cr)
	if err != nil {
		warning("error encountered updating chat room:", err.Error())
		ReportSuccess(w, false, err.Error())
		return
	}
	// report on success
	info("updated chat room:", cr.Title)
	ReportSuccess(w, true, "")
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
		ReportSuccess(w, false, err.Error())
		return
	}
	// report on success
	info("deleted chat room:", cr.Title)
	ReportSuccess(w, true, "")
	return
}
