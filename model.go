package main

import "github.com/gorilla/websocket"

type SocketPayload struct {
	Message string `json:"msg"`
	To      string `json:"to"`
	From    string `json:"from"`
}

type SocketResponse struct {
	From    string `json:"from"`
	Type    string `json:"type"`
	Message string `json:"msg"`
}

type WebSocketConnection struct {
	*websocket.Conn
	Username string
}

type Message struct {
	RoomID int
	Msg    string
}
