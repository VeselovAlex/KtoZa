package servers

import "net/http"

// Requestable is an interface for objects,
// which can be filled from request data
type Requestable interface {
	FromRequest(*http.Request)
}

type Responsable interface {
	WriteResponse(http.ResponseWriter)
}

type RequestResponsanle interface {
	Requestable
	Responsable
}
