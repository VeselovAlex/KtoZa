package main

import (
	"log"
	"net/http"
)

var addr string = ":8888"

func main() {
	pollCtrl := &PollController{}
	http.Handle("/api/poll", pollCtrl)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Unable to start master server on %s: %v\n", Addr, err)
	}
}
