package main

import "golang.org/x/net/websocket"

func handleWebSocketConnection(conn *websocket.Conn) {
	Respondents.Subscribe(conn)
}
