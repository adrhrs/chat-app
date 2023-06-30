package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	users     = make(map[string]*WebSocketConnection)
	rooms     = make(map[int][]string)
	userRooms = make(map[int][]string)
	msgRooms  = make(map[int][]Message)

	UserConnected    = "user_connected"
	UserDisconnected = "user_disconnected"
	RegularChat      = "reg_chat"
)

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "Could not open requested file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s", content)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		currentGorillaConn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
		if err != nil {
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		}

		username := r.URL.Query().Get("username")
		currentConn := WebSocketConnection{Conn: currentGorillaConn, Username: username}
		users[username] = &currentConn
		go handleIO(&currentConn, username)
	})

	log.Println("Server starting at :8080")
	http.ListenAndServe(":8080", nil)
}

func handleIO(currentConn *WebSocketConnection, username string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR", fmt.Sprintf("%v", r))
		}
	}()

	// loadPreviousMessage(currentConn, username)
	loadActiveUsers(username)
	sendMessageToAll(username, UserConnected, " joining the chat")
	// broadcastMessage(currentConn, UserJoin, "", "is joining")

	for {
		payload := SocketPayload{}
		err := currentConn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				// broadcastMessage(currentConn, UserLeft, "", "is disconnected")
				sendMessageToAll(username, UserDisconnected, " left the chat")
				delete(users, currentConn.Username)
				return
			}
		}
		sendMessageToUser(username, payload.To, RegularChat, payload.Message)
	}
}

func loadPreviousMessage(currentConn *WebSocketConnection, username string) {

}

func loadActiveUsers(username string) {
	for _, u := range users {
		if u.Username != username {
			sendMessageToUser(u.Username, username, UserConnected, " is available")
		}
	}
}

// connect and disconnect event
func sendMessageToAll(sender, chatType, message string) {
	for _, u := range users {
		if u.Username != sender {
			sendMessageToUser(sender, u.Username, chatType, message)
		}
	}
}

// chat specific func
func sendMessageToUser(sender, target, chatType, message string) {
	senderConn := users[sender]
	if targetConn, ok := users[target]; ok {
		targetConn.WriteJSON(SocketResponse{
			From:    sender,
			Type:    chatType,
			Message: message,
		})
	} else {
		senderConn.WriteJSON(SocketResponse{
			From:    "BOT",
			Type:    chatType,
			Message: "USER NOT FOUND",
		})
	}

}
