package main

import "github.com/gorilla/websocket"

type SocketPayload struct {
	Message string
	To      string
}

type SocketResponse struct {
	From    string
	Type    string
	Message string
}

type WebSocketConnection struct {
	*websocket.Conn
	Username string
}

type Message struct {
	RoomID int
	Msg    string
}
