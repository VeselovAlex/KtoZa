package main

import (
	"fmt"
	"log"
	"net/http"

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

// App реализует хранение конфигурации приложения и внедрение зависимостей
var App struct {
	Host string

	PollController       *controllers.PollController
	StatisticsController *controllers.StatisticsController
}

func init() {
	App.Host = ":8888"

	logListener := new(LogListener)

	App.StatisticsController = controllers.NewStatisticsController(logListener)
	App.PollController = controllers.NewPollController(logListener, App.StatisticsController)
}

func main() {
	fmt.Println("KtoZa poll provider. Master server")
	fmt.Println("Initialization...")

	http.Handle("/api/poll", App.PollController)
	fmt.Println("#   /api/poll")

	http.Handle("/api/stats", App.StatisticsController)
	fmt.Println("#   /api/stats")

	fmt.Println("Initialization complete")
	fmt.Println("Starting server on", App.Host)

	// Starting server on specified addr
	err := http.ListenAndServe(App.Host, nil)
	if err != nil {
		fmt.Printf("Unable to start master server on %s: %v\n", App.Host, err)
		log.Fatalf("Unable to start master server on %s: %v\n", App.Host, err)
	}
}
