# Реализация онлайн библиотека песен 🎶
# Запуск.
Для корректной работы приложения рекомендуется использовать docker и docker compose.

В начале необходимо сконфигурировать .env файл, указав все необходимые настройки. По умолчанию выставлены стандартные настройки, при которых можно из коробки проверить работоспособность.

Для создания и запуска образа в одном контейнере с PostgreSQL используйте следующую команду:
```
make run-container
```
А для запуска встроенных тестов:
```
make run-tests
```
# Проверка требований
## 1. Выставить rest методы
1. Получение данных библиотеки с фильтрацией по всем полям и пагинацией
2. Получение текста песни с пагинацией по куплетам
3. Удаление песни
4. Изменение данных песни
5. Добавление новой песни


Доступные маршруты описаны следующей функцией в файле internal/app/app.go
```go
func initRoutes(h *delivery.Handler, log *slog.Logger) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /songs", middleware.WithLogging(log, http.HandlerFunc(h.Create)))
	router.Handle("PUT /songs/{id}", middleware.WithLogging(log, http.HandlerFunc(h.Update)))
	router.Handle("GET /songs", middleware.WithLogging(log, http.HandlerFunc(h.GetAll)))
	router.Handle("GET /songs/{id}", middleware.WithLogging(log, http.HandlerFunc(h.GetVerses)))
	router.Handle("DELETE /songs/{id}", middleware.WithLogging(log, http.HandlerFunc(h.Delete)))

	// Тут находится роутер для сваггера и вывод логов

	return router
}
```
## 2. При добавлении сделать запрос в API
Для создания запросов в стороннее API был реализован отдельный клиент, расположенный в папке internal/client. Там же расположен тестовый API сервис, который можно использовать указав переменную USE_TEST_API как true в файле .env
## 3. Обогащенную информацию положить в БД postgres (структура БД должна быть создана путем миграций при старте сервиса)
Во время создания нового экземпляра приложения сервис также подключается к базе данных для применения миграций.

За это отвечает следующая функция, расположенная в файле internal/app/app.go
```go
func New(cfg *config.Config) (*App, error) {
	// Инициализируем логгер

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	m, err := migrate.New("file:///migrations", connStr)
	if err != nil {
		return nil, err
	}

	if err = m.Up(); err != nil {
		return nil, err
	}
	
	// Продолжаем создавать приложение
}
```
Во время работы приложение взаимодействует с базой данных с помощью методов, описанных в файле internal/storage/postgres/storage.go
## 4. Покрыть код Debug и Info логами
На этапе создания приложение инициализирует новый логгер с уровнем, указанным переменной LOG_LEVEL в файле .env
В процессе работы сгенерированный логгер используется на всех уровнях работы. А именно - на уровне репозитория, на уровне сервисов (там, где присутствует новая логика и где это действительно необходимо), и на уровне контроллера (хендлеров).
## 5. Вынести конфигурационные данные в .env
В файле .env расположена вся основная конфигурация приложения.
## 6. Сгенерировать сваггер на реализованное API
После создания и инициализации приложение выведет в лог на уровне Info сообщение с ссылкой на реализованный swagger. При стандартной конфигурации он будет расположен по адресу - http://localhost:8080/swagger/index.html