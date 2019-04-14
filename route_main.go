package main

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/DavidSchott/chitchat/data"
)

// GET /
// Default page
func index(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "layout", "sidebar", "public.header", "index")
}

// GET /
// Default page
func test(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "layout.1", "sidebar", "public.header", "index")
}

// GET /chat/list
func listChats(w http.ResponseWriter, r *http.Request) {
	//vals := r.URL.Query()
	rooms, err := data.CS.Chats()
	if err != nil {
		error_message(w, r, "Cannot retrieve chats")
	} else {
		p(rooms)
		generateHTML(w, &rooms, "layout", "sidebar", "public.header", "list")
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
	cr, err := data.Retrieve(title)
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(&cr, "", "\t\t")
	if err != nil {
		return
	}
	// report on success
	info(cr.User, "retrieved chat room:", cr.Title)
	p(cr)
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
	err = cr.Create()
	p("index", *data.CS.Index)
	p(data.CS.Rooms)
	// report on success/error
	if err != nil {
		warning("error encountered creating chat room:", err.Error())
		ReportSuccess(w, false, err.Error())
		return
	}
	info(cr.User, "created chat room:", cr.Title)
	ReportSuccess(w, true, "")
	return
}

// Update a room
// PUT /chat/<id>
func handlePut(w http.ResponseWriter, r *http.Request) (err error) {
	title := path.Base(r.URL.Path)
	cr, err := data.Retrieve(title)
	if err != nil {
		return
	}
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	json.Unmarshal(body, &cr)
	err = cr.Update()
	if err != nil {
		warning("error encountered updating chat room:", err.Error())
		ReportSuccess(w, false, err.Error())
		return
	}
	// report on success
	info(cr.User, "updated chat room:", cr.Title)
	ReportSuccess(w, true, "")
	return
}

// Delete a room
// DELETE /chat/<id>
func handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	title := path.Base(r.URL.Path)
	cr, err := data.Retrieve(title)
	if err != nil {
		return
	}
	err = cr.Delete()
	if err != nil {
		warning("error encountered deleting chat room:", err.Error())
		ReportSuccess(w, false, err.Error())
		return
	}
	// report on success
	info(cr.User, "deleted chat room:", cr.Title)
	ReportSuccess(w, true, "")
	return
}
