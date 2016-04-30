package main

import (
	"encoding/json"
	"io"
)

var JSONEnqDeq *jsonEnqDec

type jsonEnqDec struct{}

func (*jsonEnqDec) FromRequest(r io.Reader, to interface{}) error {
	return json.NewDecoder(r).Decode(to)
}

func (*jsonEnqDec) ToResponseWriter(w io.Writer, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
