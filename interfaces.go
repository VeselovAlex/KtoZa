package KtoZa

import (
	"io"

	"github.com/VeselovAlex/KtoZa/model"
)

// PollStorage представляет стандартный интерфейс хранилища
// опросов в системе KtoZa
type PollStorage interface {
	// Возвращает текущий опрос или nil, если опрос не задан
	Get() *model.Poll

	// Создает новый или обновляет текущий опрос. Возвращает
	// указатель на обновленный опрос или nil, если опрос
	// не был обновлен
	CreateOrUpdate(poll *model.Poll) *model.Poll

	// Удаляет текущий опрос. Возвращает указатель на удаленный
	// опрос или nil, если опрос был удален ранее
	Delete() *model.Poll
}

// StatisticsStorage представляет стандартный интерфейс хранилища
// статистики в системе KtoZa
type StatisticsStorage interface {
	// Возвращает указатель на текущую статистику или nil,
	// если статистика не создана
	Get() *model.Statistics

	// Создает новую или объединяет текущую статистику с заданной. Возвращает
	// указатель на обновленную статистику или nil, если статистика
	// не был обновлена
	CreateOrJoinWith(*model.Statistics) *model.Statistics

	// Удаляет текущую статистику. Возвращает указатель на удаленную
	// статистику или nil, если статистика была удалена ранее
	Delete() *model.Statistics
}

// PubSub представляет интерфейс контроллера соединений модели Publisher/Subcriptor
type PubSub interface {
	Subscribe(io.ReadWriteCloser)
	NotifyAll(interface{})
	Await() interface{}
}
