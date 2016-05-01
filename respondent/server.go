package main

import (
	"log"
	"net/http"

	"github.com/VeselovAlex/KtoZa/model"
	"github.com/VeselovAlex/KtoZa/respondent/controllers"
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

var App struct {
	Name    string
	Version string
	Host    string

	PollController       *controllers.PollController
	StatisticsController *controllers.StatisticsController
	AnswerController     *controllers.AnswerController
	SessionController    *controllers.SessionController
}

func init() {
	// Configure App
	App.Name = "KtoZa poll provider, respondent server"
	App.Version = "0.0.1"
	App.Host = ":8080"

	logListener := &LogListener{}

	App.StatisticsController = controllers.NewTestStatCtrl(logListener)
	App.SessionController = controllers.NewSessionController()
	App.PollController = controllers.NewTestPollCtrl(logListener, App.StatisticsController, App.SessionController)
	App.AnswerController = controllers.NewTestAnswerCtrl(logListener, App.StatisticsController)
}

func main() {
	//statics := http.FileServer(http.Dir("client"))
	//http.Handle("/", http.StripPrefix("/", statics))
	http.Handle("/api/poll", App.PollController)
	http.Handle("/api/stats", App.StatisticsController)
	http.Handle("/api/register", App.SessionController)
	http.Handle("/api/submit", App.AnswerController)
	log.Println("Starting server on", App.Host)
	err := http.ListenAndServe(App.Host, nil)
	if err != nil {
		log.Fatalln("Server stopped cause:", err)
	}
}
