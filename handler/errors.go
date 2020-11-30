package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/DavidSchott/chitchat/data"
)

type errHandler func(http.ResponseWriter, *http.Request) error

func (fn errHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		if apierr, ok := err.(*data.APIError); ok {
			w.Header().Set("Content-Type", "application/json")
			apierr.SetMsg()
			warning("API error:", apierr.Error())
			if apierr.Code == 101 || apierr.Code == 201 {
				notFound(w, r)
			}
			if apierr.Code == 102 || apierr.Code == 202 || apierr.Code == 303 {
				badRequest(w, r)
			}
			ReportSuccess(w, false, apierr)
		} else {
			danger("Server error", err.Error())
			http.Error(w, err.Error(), 500)
		}
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	info("Not found request:", r.RequestURI)
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	info("Bad request:", r.RequestURI, r.Body)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "/index.html")
	w.WriteHeader(302)
}

// Convenience function to redirect to the error message page
func errorMessage(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(writer, request, strings.Join(url, ""), 302)
}

// GET /err?msg=
// shows the error message page
func handleError(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	fmt.Fprintf(writer, "Error: %s!", vals.Get("msg"))
	warning(fmt.Sprintf("Error: %s!", vals.Get("msg")))
}
