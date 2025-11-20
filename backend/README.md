# Backend module — Orders (Go)

Документация по backend-модулю проекта Orders.

## Краткое описание

Backend реализован на Go с использованием фреймворка Gin для HTTP и GORM для работы с PostgreSQL. Модуль содержит модели, репозитории, сервисный слой и HTTP-обработчики (handlers), которые обрабатывают REST API для управления пользователями, клиентами, договорами и заказами.

## Структура (важные файлы)

- `cmd/server/main.go` — точка входа: загрузка конфигурации, подключение к БД, миграции и запуск роутера.
- `internal/config/config.go` — конфигурация приложения (DSN, секреты JWT и пр.).
- `internal/models/models.go` — модели GORM: `User`, `Client`, `Contract`, `ContractAddress`, `Product`, `Order`, `OrderItem` и др.
- `internal/repository` — слои репозиториев для прямой работы с GORM (CRUD операции).
- `internal/service` — бизнес-логика: обёртка над репозиториями, валидация, дополнительные операции.
- `internal/api/handlers.go` — регистрация маршрутов и общие HTTP-обработчики (login, signup, защищённые маршруты).
- `internal/api/clients.go` *(рекомендуется)* — отдельные handlers для работы с клиентами/договорами (если добавлены).

> Для поиска конкретных точек входа используйте `grep`/IDE: `SetupRoutes`, `NewService`, `AutoMigrate`.

## Архитектурные принципы

- Чёткое разделение ответственности: handlers → service → repository → models.
- Авторизация через JWT: `/login` возвращает токен, который используется в заголовке `Authorization: Bearer <token>` для защищённых маршрутов.
- Миграции выполняются автоматически при старте (`AutoMigrate`) — проверяйте реальную схему при продакшн-развёртывании.

## Модели (коротко)

- User — пользователь, хранит email, пароль (хеш), роль и ссылку на канал продаж (CanalID / Canal). Может быть nullable.
- Client — клиент (заказчик): поля `Name`, `Email`, `Phone`, `Address`, `UserID` (владелец).
- Contract — договор клиента: `ClientID`, `Number`, `Date`, `Amount`, `Status` и связи на `Client` и `ContractAddress`.
- ContractAddress, Product, Order, OrderItem — вспомогательные сущности для адресов, товаров и заказов.

Если нужно, откройте `internal/models/models.go` для полной структуры и тегов GORM.

## Основные API endpoints (пример)

Все примеры предполагают, что сервер слушает на `http://localhost:8080`.

- POST /login — вход: принимает JSON `{ "email": "...", "password": "..." }`, возвращает JSON `{ "token": "..." }`.

- POST /clients — создать клиента (защищённый): header `Authorization: Bearer <token>`; body: `{ "name":"ACME", "email":"acme@example.com", "phone":"...", "address":"..." }`. `UserID` рекомендуется брать из токена на стороне сервера.

- GET /clients/:id — получить клиента по id (включая связанные договоры, если service делает Preload).

- POST /clients/:id/contracts — создать договор для клиента: body `{ "number":"CTR-001","date":"2025-11-01","amount":1000.0,"status":"active" }`.

- GET /contracts/:id — получить договор по id.

Проверьте `internal/api/handlers.go` для полной карты маршрутов и middlewares.

## Примеры curl

Получить токен (login):

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret"}'
```

Создать клиента (замените `<token>`):

```bash
curl -X POST http://localhost:8080/clients \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"ACME","email":"acme@example.com","phone":"123","address":"addr"}'
```

Создать договор для клиента id=1:

```bash
curl -X POST http://localhost:8080/clients/1/contracts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"number":"CTR-001","date":"2025-11-01","amount":1000.00,"status":"active"}'
```

## JWT и тестирование в Postman

1. Выполните `POST /login` с реальными credentials, сохраните `token` в окружении Postman (Tests script):

```javascript
const json = pm.response.json();
pm.environment.set("auth_token", json.token);
```

2. В Authorization для последующих запросов выберите `Bearer Token` и используйте `{{auth_token}}`.

3. Если нужно вручную сгенерировать JWT в Pre-request, используйте Secret и алгоритм (HS256). Но безопаснее получать токен от сервера.

## Миграции и развёртывание

- В `cmd/server/main.go` выполняется `db.AutoMigrate(...)` при старте. Для продакшна лучше использовать явные миграции (migrate tool) и процессы отката.
- Конфигурация БД берётся из `internal/config` (DSN). Убедитесь, что переменные окружения заданы.

## Инициализация сервиса (пример)

В `cmd/server/main.go` создаётся репозиторий и сервис и передаётся в маршруты:

```go
repo := repository.NewRepository(db)
svc := service.NewService(repo, cfg.JWTSecret)
api.SetupRoutes(r, svc)
```

Service содержит методы для: регистрации/логина, создания клиентов/договоров, получения сущностей и т.д.

## Советы по безопасности и валидации

- Пароли храните только в виде хеша (bcrypt).
- Проверяйте ownership: client.UserID следует заполнять серверно (из токена), а не позволять клиенту присылать чужой user_id.
- Валидация уникальности (email, номер договора) — проверяйте через DB и возвращайте понятные ошибки.

## Что можно улучшить / next steps

- Добавить Postman collection с примером логина и CRUD для клиентов/договоров.
- Добавить unit/integration тесты для сервисного слоя и handlers (httptest + gin).
- Перевести миграции на инструмент (golang-migrate) для управляемых миграций.
- Добавить логирование запросов и ошибок (structured logging).

## Контакты в коде

- API-обработчики: `internal/api/handlers.go` (+ дополнительные файлы в `internal/api/`).
- Бизнес-логика: `internal/service/*`.
- Репозитории/DB: `internal/repository/*`.
- Модели: `internal/models/models.go`.

---

Если хотите — могу:
- добавить Postman collection и пример environment variables;
- сделать `clients.go` с полным CRUD и примерами тестов;
- подготовить пример миграции через `golang-migrate`.
