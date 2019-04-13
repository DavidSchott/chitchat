package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/DavidSchott/chitchat/data"
)

// GET /
// Default page
func index(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "layout", "public.header", "index")
}

// POST /chat
// Create chat page
func create(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	classification, err := strconv.ParseInt(r.MultipartForm.Value["classification"][0], 10, 0)
	if err != nil {
		p(err)
		danger("Could not create chat page", err)
		fmt.Fprintf(w, "Error creating chat")
	} else {
		cr := &data.ChatRoom{
			Title:    r.MultipartForm.Value["title"][0],
			User:     r.MultipartForm.Value["name"][0],
			Type:     uint(classification),
			Password: r.MultipartForm.Value["secret"][0],
			ID:       chatServer.Index,
		}
		chatServer.Push(cr)
		chatServer.Index++
		generateHTML(w, cr.Type == 0, "layout", "public.header", "chat")
		info(cr.User, "created chat room:", cr.Title)
	}
}

// User joins a chat room
func join(w http.ResponseWriter, r *http.Request) {
}

// User deletes a chat room
func delete(w http.ResponseWriter, r *http.Request) {
}

// Check for duplicate chatrooms
func check(w http.ResponseWriter, r *http.Request) {
}

// List chatrooms
func list(w http.ResponseWriter, r *http.Request) {
}
