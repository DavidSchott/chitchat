package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Post struct {
	User string
	Msg  string
	Time string
}

func jsonExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	now := time.Now()
	var time string = fmt.Sprintf("%d:%d", now.Hour(), now.Minute())
	fmt.Println(time)
	post := &Post{
		User: "David Schott",
		Msg:  "hello world",
		Time: time,
	}
	json, _ := json.Marshal(post)
	w.Write(json)
}
