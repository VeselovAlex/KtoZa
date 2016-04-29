package main

<<<<<<< HEAD
import common "github.com/VeselovAlex/KtoZa"

// App представляет конфигурацию приложения
var App common.Config
=======
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
>>>>>>> 6e99eb4... session controller, app configuration added
