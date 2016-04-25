package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"time"

	ws "golang.org/x/net/websocket"
)

var serverAddr = new(string)
var templateDir = new(string)

func init() {
	*serverAddr = ":8080"
	*templateDir = "templates"
}

type templatehandler struct {
	TemplateName string
	Data         interface{}
	t            *template.Template
}

func (h *templatehandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.t = template.Must(template.ParseFiles(path.Join(*templateDir, h.TemplateName)))
	h.t.Execute(w, h.Data)
}

func getPollInfo(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(PollController.Get()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func inTime() bool {
	poll := PollController.Get()
	now := time.Now()
	return poll.StartsAt.Before(now) && poll.EndsAt.After(now)
}

func main() {
	Master.Addr = "http://localhost:8888/api/"
	PollController.Init()
	go StatisticsCache.Run()
	go PollController.Run()
	hello := &templatehandler{TemplateName: "hello.html"}
	http.Handle("/", hello)
	http.HandleFunc("/poll", getPollInfo)
	http.HandleFunc("/submit", submitUserAnswer)

	http.Handle("/ws", ws.Handler(handleWSConnection))
	http.ListenAndServe(*serverAddr, nil)
}
