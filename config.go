package KtoZa

type Config struct {
	// Информация о приложении

	Name    string
	Version string

	// Конфигурация запуска
	Host string

	// Конфигурация хранилищ
	PollStorage       PollStorage
	StatisticsStorage StatisticsStorage

	// Конфигурация сериализаторов
	RequestDecoder  RequestDecoder
	ResponseEncoder ResponseEncoder
}
