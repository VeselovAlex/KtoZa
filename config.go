package KtoZa

// Config представляет конфигурацию сервера системы KtoZa
type Config struct {
	// Информация о приложении
	Name    string
	Version string

	// Конфигурация запуска
	Host string

	// Конфигурация хранилищ
	PollStorage       PollStorage
	StatisticsStorage StatisticsStorage

	// Конфигурация Pub/Sub
	PubSub PubSub
}
