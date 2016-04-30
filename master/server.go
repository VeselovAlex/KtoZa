package main

import (
	"fmt"
	"log"
	"net/http"

	common "github.com/VeselovAlex/KtoZa"
)

var addr string

// App реализует хранение конфигурации приложения и внедрение зависимостей
var App common.Config

func init() {
	App.Host = ":8888"

	App.PollStorage = NewMasterPollStorage()
	App.StatisticsStorage = NewMasterStatisticsStorage()

	App.PubSub = newWebSocketPubSub()
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

	// WebSocketPubSub starts
	http.Handle("/ws", NewWebSocketPubSubController())
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
