package KtoZa

import (
	"io"

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
	FromRequest(io.Reader, interface{}) error
}

type ResponseEncoder interface {
	ToResponseWriter(io.Writer, interface{}) error
}
