package KtoZa

import "net/http"

type Answer []answerOptions

type answerOptions []string

func (a *Answer) FromRequest(r *http.Request) {
}
