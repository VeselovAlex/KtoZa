package main

import (
	core "KtoZa"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

var Master masterServer

type masterServer struct {
	Addr string
}

func (m masterServer) GetPoll() *core.Poll {
	resp, err := http.Get(m.Addr + "poll")
	if err != nil {
		log.Println("Error while requesting poll:", err)
		return nil
	}
	poll := &core.Poll{}
	err = json.NewDecoder(resp.Body).Decode(poll)
	if err != nil {
		log.Println("Error while parsing poll:", err)
		return nil
	}
	return poll
}

func (m masterServer) SendSubmit(stat *core.Statistics) error {
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	err := json.NewEncoder(buffer).Encode(stat)
	if err != nil {
		return err
	}
	_, err = http.Post(m.Addr+"submit", "application/json", buffer)
	return err
}
