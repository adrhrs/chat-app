package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	users = make(map[string]*WebSocketConnection)
	// rooms     = make(map[int][]string)
	// userRooms = make(map[int][]string)
	// msgRooms  = make(map[int][]Message)

	// host             = "0.0.0.0"
	listUser    = "list_user"
	RegularChat = "reg_chat"

	botUser = "bot"

	port = "8080"
	host = "127.0.0.1"
)

func main() {
	http.HandleFunc("/ping", ping)

	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/ws", handleWebSocket)

	log.Println(fmt.Sprintf("Chat App starting at %v:%v", host, port))
	http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	currentGorillaConn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	username := r.URL.Query().Get("username")
	currentConn := WebSocketConnection{Conn: currentGorillaConn, Username: username}
	users[username] = &currentConn
	go handleIO(&currentConn, username)
}

func handleIO(currentConn *WebSocketConnection, username string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR", fmt.Sprintf("%v", r))
		}
	}()

	loadActiveUsers(username)
	log.Println(username, "is connected", users)

	for {
		payload := SocketPayload{}
		err := currentConn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				delete(users, currentConn.Username)
				loadActiveUsers(username)
				log.Println(username, "is disconnected", users)
				return
			}
		}
		fmt.Println(payload)
		sendMessageToUser(username, payload.To, RegularChat, payload.Message)
	}
}

func loadActiveUsers(username string) {
	activeUsers := []string{}
	for username := range users {
		activeUsers = append(activeUsers, username)
	}

	message := (strings.Join(activeUsers[:], ","))
	sendMessageToAll(botUser, listUser, message)

}

func sendMessageToAll(sender, chatType, message string) {
	for _, u := range users {
		sendMessageToUser(sender, u.Username, chatType, message)
	}
}

func sendMessageToUser(sender, target, chatType, message string) {
	if targetConn, ok := users[target]; ok {
		targetConn.WriteJSON(SocketResponse{
			From:    sender,
			Type:    chatType,
			Message: message,
		})
	}
}
