package main

import (
	"log"
	"net/http"
	"os"

	"github.com/VeselovAlex/KtoZa/model"

	"github.com/VeselovAlex/KtoZa/master/controllers"
)

type LogListener struct{}

func (l *LogListener) OnPollUpdate(*model.Poll) {
	log.Println("LOG LISTENER :: Poll update")
}

func (l *LogListener) OnStatisticsUpdate(*model.Statistics) {
	log.Println("LOG LISTENER :: Statistics update")
}

func (l *LogListener) OnNewAnswerSet(ans model.AnswerSet) {
	log.Println("LOG LISTENER :: New answer set", ans)
}

var appHost string

func init() {
	appHost = ":8888"

	log.Println("SERVER INIT :: KtoZa poll provider. Master server")
	log.Println("SERVER INIT :: Initialization started...")

	log.Println("SERVER INIT :: Opening server data folder...")
	data := "data"
	// Создаем папку выходных данных
	err := os.Mkdir(data, 0755)
	if err != nil && !os.IsExist(err) {
		// Папка не создана и не существует
		log.Fatalln("SERVER INIT :: Initialization failed:", err)
	}

	log.Println("SERVER INIT :: Loading data storage...")
	controllers.LoadFileSystemStorage(data)

	log.Println("SERVER INIT :: Initializing request handlers")

	log.Println("SERVER INIT :: #   /api/ws")
	wsCtrl := controllers.NewWebSocketController()
	http.Handle("/api/ws", wsCtrl)

	log.Println("SERVER INIT :: #   /api/stats")
	statCtrl := controllers.NewStatisticsController(wsCtrl)
	http.Handle("/api/stats", statCtrl)

	log.Println("SERVER INIT :: #   /api/poll")
	pollCtrl := controllers.NewPollController(statCtrl, wsCtrl)
	http.Handle("/api/poll", pollCtrl)

	log.Println("SERVER INIT :: Initialization complete")
}

func main() {
	log.Println("SERVER INIT :: Starting server on", appHost)
	// Starting server on specified addr
	err := http.ListenAndServe(appHost, nil)
	if err != nil {
		log.Fatalf("SERVER INIT :: Unable to start master server on %s: %v\n", appHost, err)
	}
}
