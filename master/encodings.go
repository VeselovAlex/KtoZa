package main

import (
	"encoding/json"
	"net/http"
)

var JSONEnqDeq *jsonEnqDec

type jsonEnqDec struct{}

func (*jsonEnqDec) FromRequest(r *http.Request, to interface{}) error {
	return json.NewDecoder(r.Body).Decode(to)
}

func (*jsonEnqDec) ToResponseWriter(w http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
