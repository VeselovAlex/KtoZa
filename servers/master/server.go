package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"

	core "KtoZa"
)

var templateDir = "templates"

func main() {
	poll := core.GetDummyPoll()
	PollStorage.Create(poll)
	StatisticsCtrl.Init(core.NewStatistics(poll))
	http.HandleFunc("/api/poll", handleGetPoll)
	http.HandleFunc("/api/stats", handleGetStats)
	http.HandleFunc("/api/submit", handleSubmitAnswerCache)
	http.ListenAndServe(":8888", nil)
}

// Use JSON Data
func handleNewPoll(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	poll := &core.Poll{}
	if err := decoder.Decode(poll); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Add error handling
	PollStorage.Create(poll)
}

func handleNewPollPage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(
		template.ParseFiles(
			path.Join(templateDir, "create.html")))

	t.Execute(w, nil)
}

func handleNewPollRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleNewPollPage(w, r)
	case "POST":
		handleNewPoll(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSubmitAnswerCache(w http.ResponseWriter, r *http.Request) {
	log.Println("Submission request")
	stats := &core.Statistics{}
	err := json.NewDecoder(r.Body).Decode(stats)
	if err != nil {
		http.Error(w, "Bad answer cache format", http.StatusBadRequest)
		return
	}
	StatisticsCtrl.JoinWith(stats)
}

func handleGetPoll(w http.ResponseWriter, r *http.Request) {
	poll, err := PollStorage.Get()
	if err != nil {
		http.Error(w, "No poll", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(poll)
}

func handleGetStats(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(StatisticsCtrl.Get())
}
