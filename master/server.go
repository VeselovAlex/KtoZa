package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VeselovAlex/KtoZa/master/controllers"
)

var appHost string
var pwd string
var needUsage bool
var logFile string

func init() {
	// Чтение параметров командной строки
	flag.StringVar(&appHost, "a", ":8888", "This server host")
	flag.BoolVar(&needUsage, "h", false, "Show this help")
	flag.StringVar(&logFile, "f", "", "Write log to file")
	flag.StringVar(&pwd, "p", "", "Master server password, required")
	flag.Parse()

	if needUsage || pwd == "" {
		showUsage()
		os.Exit(-1)
	}

}

func main() {
	fmt.Println("KtoZa poll provider. Master server")

	// Загрузка лог-файла
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_APPEND, 0755)
		if os.IsNotExist(err) {
			file, err = os.Create(logFile)
		}
		if err != nil {
			fmt.Println("Can not open log file:", err)
			os.Exit(1)
		}
		defer file.Close()
		log.SetOutput(file)
	}

	log.Println("SERVER INIT :: Initialization started...")

	loadDataStorage()
	initHandlers()

	log.Println("SERVER INIT :: Starting server on", appHost)
	err := http.ListenAndServe(appHost, nil)
	if err != nil {
		log.Fatalln("SERVER INIT :: Server stoppped cause:", err)
	}
}

func initHandlers() {
	log.Println("SERVER INIT :: Initializing request handlers")

	log.Println("SERVER INIT :: #   /api/auth")
	auth := controllers.NewUserControl(pwd)
	http.Handle("/api/auth", auth)

	log.Println("SERVER INIT :: #   /api/ws")
	wsCtrl := controllers.NewWebSocketController()
	http.Handle("/api/ws", wsCtrl)

	log.Println("SERVER INIT :: #   /api/stats")
	statCtrl := controllers.NewStatisticsController(wsCtrl)
	http.Handle("/api/stats", auth.Authorized(statCtrl, "DELETE"))

	log.Println("SERVER INIT :: #   /api/poll")
	pollCtrl := controllers.NewPollController(wsCtrl, statCtrl)
	http.Handle("/api/poll", auth.Authorized(pollCtrl, "PUT", "DELETE"))

	log.Println("SERVER INIT :: #   /")
	fs := http.FileServer(http.Dir("client"))
	http.Handle("/", http.StripPrefix("/", fs))

	log.Println("SERVER INIT :: Initialization complete")
}

func loadDataStorage() {

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
}

func showUsage() {
	flag.Usage()
}
