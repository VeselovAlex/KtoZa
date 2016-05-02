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

var appHost string
var masterHost string

func init() {
	appHost = ":8080"
	masterHost = "http://localhost:8888"

	log.Println("SERVER INIT :: KtoZa poll provider. Respondent server")
	log.Println("SERVER INIT :: Initialization started...")

	log.Println("SERVER INIT :: Connecting to master...")
	controllers.ConnectToMaster(masterHost)
	log.Println("SERVER INIT :: Connected to master")

	log.Println("SERVER INIT :: Initializing request handlers")

	log.Println("SERVER INIT :: #   /api/stats")
	statCtrl := controllers.NewStatisticsController()
	http.Handle("/api/stats", statCtrl)

	log.Println("SERVER INIT :: #   /api/register")
	sessionCtrl := controllers.NewSessionController()
	http.Handle("/api/register", sessionCtrl)

	log.Println("SERVER INIT :: #   /api/submit")
	ansCtrl := controllers.NewAnswerController(statCtrl)
	http.Handle("/api/submit", ansCtrl)

	log.Println("SERVER INIT :: #   /api/poll")
	pollCtrl := controllers.NewPollController(statCtrl, sessionCtrl, ansCtrl)
	http.Handle("/api/poll", pollCtrl)

	log.Println("SERVER INIT :: Initialization complete")
}

func main() {
	log.Println("SERVER INIT :: Starting server on", appHost)
	err := http.ListenAndServe(appHost, nil)
	if err != nil {
		log.Fatalln("SERVER INIT :: Server stopped cause:", err)
	}
}
