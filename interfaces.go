package KtoZa

import (
	"net/http"

	"github.com/VeselovAlex/KtoZa/model"
)

type PollStorage interface {
	Get() *model.Poll
	CreateOrUpdate(poll *model.Poll) *model.Poll
	Delete() *model.Poll
}

type StatisticsStorage interface {
	Get() *model.Statistics
	CreateOrUpdate(*model.Statistics) *model.Statistics
	Delete() *model.Statistics
}

type RequestDecoder interface {
	FromRequest(*http.Request, interface{}) error
}

type ResponseEncoder interface {
	ToResponseWriter(http.ResponseWriter, interface{}) error
}
