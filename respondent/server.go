package main

import (
	"log"
	"net/http"
)

var addr string

func init() {
	// Configure App
	App.Name = "KtoZa poll provider, respondent server"
	App.Version = "0.0.1"
	App.Host = ":8080"
}

func main() {
	//statics := http.FileServer(http.Dir("client"))
	//http.Handle("/", http.StripPrefix("/", statics))
	http.Handle("/api/poll", &PollController{})
	http.Handle("/api/stats", &StatisticsController{})
	//http.Handle("/api/register", &SessionController{})
	http.Handle("/api/submit", &AnswerController{})
	log.Println("Starting server on", App.Host)
	err := http.ListenAndServe(App.Host, nil)
	if err != nil {
		log.Fatalln("Server stopped cause:", err)
	}
}
