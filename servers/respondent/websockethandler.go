package main

import (
	ws "golang.org/x/net/websocket"
)

func handleWSConnection(conn *ws.Conn) {
	PollController.SendSecretAndClose(conn)
}
