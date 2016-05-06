package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VeselovAlex/KtoZa/master/controllers"
)

type Authorized struct {
	hash string
}

func NewAuthorized(hash string) *Authorized {
	return &Authorized{hash}
}

func (a *Authorized) New(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

var appHost string

func init() {
	appHost = ":8888"
}

func main() {
	fmt.Println(" KtoZa poll provider. Master server")
	var err error
	/*
		hash, err := ioutil.ReadFile("hash.dat")
		if err != nil {
			got := false
			var pwd string
			var conf string
			for !got {
				fmt.Println("Введите новый пароль:")
				fmt.Scanln(&pwd)
				fmt.Println("Подтвердите пароль:")
				fmt.Scanln(&conf)
				got = conf == pwd && pwd != ""
			}
			h := md5.New()
			io.WriteString(h, pwd)
			hash = h.Sum(nil)
			ioutil.WriteFile("hash.dat", hash, 0755)
		}
		auth := NewAuthorized(string(hash))
	*/
	log.Println("SERVER INIT :: Initialization started...")
	log.Println("SERVER INIT :: Opening server data folder...")
	data := "data"
	// Создаем папку выходных данных
	err = os.Mkdir(data, 0755)
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
	pollCtrl := controllers.NewPollController(wsCtrl, statCtrl)
	http.Handle("/api/poll", pollCtrl)

	log.Println("SERVER INIT :: #   /")
	fs := http.FileServer(http.Dir("client"))
	http.Handle("/", http.StripPrefix("/", fs))

	log.Println("SERVER INIT :: Initialization complete")
	log.Println("SERVER INIT :: Starting server on", appHost)
	// Starting server on specified addr
	err = http.ListenAndServe(appHost, nil)
	if err != nil {
		log.Fatalf("SERVER INIT :: Unable to start master server on %s: %v\n", appHost, err)
	}
}
