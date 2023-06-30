package main

import (
	"encoding/json"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {

	activeUsers := []string{}
	for k := range users {
		activeUsers = append(activeUsers, k)
	}

	type Resp struct {
		ActiveUsers []string
		Msg         map[string][]string
	}

	resp := Resp{
		ActiveUsers: activeUsers,
	}
	json.NewEncoder(w).Encode(resp)

}
