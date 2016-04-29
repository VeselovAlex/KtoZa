package main

// App представляет конфигурацию приложения
var App respondentConfig

type appConfig struct {
	// Информация о приложении

	Name    string
	Version string

	// Конфигурация запуска
	Host string
}

type respondentConfig struct {
	// Базовая конфигурация
	appConfig

	// Конфигурация хранилищ
	PollStorage       PollStorage
	StatisticsStorage StatisticsStorage
}
