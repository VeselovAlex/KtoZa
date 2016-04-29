package main

import (
	"fmt"
	"log"
	"net/http"
)

var addr string

func init() {
	addr = ":8888"
}

func main() {
	fmt.Println("KtoZa poll provider. Master server")
	fmt.Println("Initialization...")

	// PollController starts
	pollCtrl := &PollController{}
	http.Handle("/api/poll", pollCtrl)
	fmt.Println("#   /api/poll")

	// StatisticsController starts
	statsCtrl := &StatisticsController{}
	http.Handle("/api/stats", statsCtrl)
	fmt.Println("#   /api/stats")

	fmt.Println("Initialzation complete")
	fmt.Println("Starting server on", addr)

	// Starting server on specified addr
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("Unable to start master server on %s: %v\n", addr, err)
		log.Fatalf("Unable to start master server on %s: %v\n", addr, err)
	}
}
