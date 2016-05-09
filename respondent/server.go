package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VeselovAlex/KtoZa/respondent/controllers"
)

var appHost string
var masterHost string
var needUsage bool
var logFile string

func init() {
	flag.StringVar(&appHost, "a", ":8080", "This server host")
	flag.StringVar(&masterHost, "m", "", "Master server address, required")
	flag.BoolVar(&needUsage, "h", false, "Show this help")
	flag.StringVar(&logFile, "f", "", "Write log to file")

	flag.Parse()

	if needUsage || masterHost == "" {
		showUsage()
		os.Exit(1)
	}

}

func main() {

	if logFile != "" {
		l, err := os.OpenFile(logFile, os.O_APPEND, 0755)
		if err != nil {
			if os.IsNotExist(err) {
				l, err = os.Create(logFile)
			}
			if err != nil {
				fmt.Println("Unable to open log file:", err)
				os.Exit(1)
			}
		}
		defer l.Close()
		log.SetOutput(l)
	}

	log.Println("SERVER INIT :: KtoZa poll provider. Respondent server")
	log.Println("SERVER INIT :: Initialization started...")

	log.Println("SERVER INIT :: Connecting to master...")
	controllers.ConnectToMaster(masterHost, "http://localhost"+appHost)
	log.Println("SERVER INIT :: Connected to master")

	initHandlers()

	log.Println("SERVER INIT :: Starting server on", appHost)
	err := http.ListenAndServe(appHost, nil)
	if err != nil {
		log.Fatalln("SERVER INIT :: Server stopped cause:", err)
	}
}

func initHandlers() {
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

	log.Println("SERVER INIT :: #   /")
	fs := http.FileServer(http.Dir("client"))
	http.Handle("/", http.StripPrefix("/", fs))

	log.Println("SERVER INIT :: Initialization complete")
}

func showUsage() {
	flag.Usage()
}
