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
			//			json, _ := json.Marshal(apierr)
			//			w.Write(json)
			warning("API error:", apierr.Error())
			ReportSuccess(w, false, err.(*data.APIError))
		} else {
			danger("Server error", err.Error())
			http.Error(w, err.Error(), 500)
		}
	}
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	fmt.Fprintln(w, "No such service, try next door")
}

func redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "http://google.com")
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
