# Backend — Development setup

Коротко: в этом файле описаны шаги для локальной разработки backend-а проекта Orders.

## Предварительные требования

- Go 1.20+ (установите с https://go.dev)
- PostgreSQL (локально или через Docker)
- git
- (опционально) Docker и Docker Compose для быстрого развёртывания БД

## Быстрый запуск с Docker (Postgres)

1. Перейдите в папку `backend`:

```powershell
cd d:\Documents\Projects\9_Golang\Orders\backend
```

2. Создайте файл `.env` на основе `.env.example` и при необходимости измените DSN/JWT_SECRET:

```powershell
cp .env.example .env
# отредактируйте .env в любом редакторе
```

3. Запустите Postgres через Docker Compose (если у вас есть `docker-compose.yml` в этой папке):

```powershell
docker compose up -d
```

4. Проверьте, что Postgres доступен, создайте базу данных, если нужно (примеры для psql):

```powershell
psql "host=localhost port=5432 user=postgres password=postgres" -c "CREATE DATABASE orders_dev;"
```

## Локальный запуск (Go)

1. Убедитесь, что в `.env` корректный `DSN` и `JWT_SECRET`.
2. Установите зависимости и запустите сервер:

```powershell
# из папки backend
go mod tidy
go run ./cmd/server
```

3. При старте `main.go` выполняется `db.AutoMigrate(...)` — это создаст необходимые таблицы автоматически.

## Как работать с миграциями

Автоматические `AutoMigrate` хорошо подходят для разработки. Для production рекомендуется использовать инструмент миграций (например, golang-migrate) и хранить SQL миграции в репозитории.

## Переменные окружения и конфигурация

- `DSN` — строка подключения к Postgres, см. `.env.example`.
- `JWT_SECRET` — секрет для подписи JWT.
- `PORT` — порт для сервера (по умолчанию 8080).

В коде конфигурация загружается через `internal/config`. Убедитесь, что при старте приложения эти переменные доступны (например, через shell или `.env` менеджер).

## Тестирование API

1. Получите токен:

```powershell
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"email":"user@example.com","password":"secret"}'
```

2. Используйте токен в `Authorization: Bearer <token>` для защищённых вызовов (например, создание клиента).

## Полезные команды

- Остановить контейнеры: `docker compose down`
- Просмотреть логи сервера: `go run ./cmd/server` (stdout) или если запущено в фоне — смотреть systemctl/pm2/compose logs

## Советы

- Не храните реальные секреты в репозитории. Используйте `.env` и git-ignore.
- В разработке ставьте понятный, но простой `JWT_SECRET`, в проде используйте безопасный секрет.
