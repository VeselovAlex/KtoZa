package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"

	common "github.com/VeselovAlex/KtoZa"
)

var addr string

// App реализует хранение конфигурации приложения и внедрение зависимостей
var App common.Config

func init() {
	App.Host = ":8888"

	App.PollStorage = NewDummyPollStorage()
	App.ResponseEncoder = JSONEnqDeq
	App.RequestDecoder = JSONEnqDeq
}

func main() {
	fmt.Println("KtoZa poll provider. Master server")
	fmt.Println("Initialization...")

	// PollController starts
	pollCtrl := &PollController{}
	http.Handle("/api/poll", pollCtrl)
	fmt.Println("#   /api/poll")

	// StatisticsController starts
	statsCtrl := &StatisticsController{}
	http.Handle("/api/stats", statsCtrl)
	fmt.Println("#   /api/stats")

	// StatisticsController starts
	wsCtrl := websocket.Handler(handleWebSocketConnection)
	http.Handle("/ws", wsCtrl)
	fmt.Println("#   /ws")

	fmt.Println("Initialization complete")
	fmt.Println("Starting server on", App.Host)

	// Starting server on specified addr
	err := http.ListenAndServe(App.Host, nil)
	if err != nil {
		fmt.Printf("Unable to start master server on %s: %v\n", App.Host, err)
		log.Fatalf("Unable to start master server on %s: %v\n", App.Host, err)
	}
}
