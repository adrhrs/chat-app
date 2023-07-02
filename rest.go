package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ping(w http.ResponseWriter, req *http.Request) {

	activeUsers := []string{}
	for k := range users {
		activeUsers = append(activeUsers, k)
	}

	type Resp struct {
		ActiveUsers []string
		Msg         string
	}

	resp := Resp{
		ActiveUsers: activeUsers,
		Msg:         "pong",
	}
	json.NewEncoder(w).Encode(resp)

}

func serveHTML(w http.ResponseWriter, req *http.Request) {

	content, err := ioutil.ReadFile("static/chat-v2.html")
	if err != nil {
		http.Error(w, "Could not open requested file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", content)

}
