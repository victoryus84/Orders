# Orders Project Structure

## Основные папки

- **frontend/** — клиентская часть на Flutter
  - **lib/pages/** — страницы приложения (экран логина, меню, создание заказа и др.)
  - **lib/models/** — модели данных (например, User, Order)
  - **lib/services/** — сервисы для работы с API и авторизацией
  - **lib/widgets/** — переиспользуемые виджеты (кнопки, выпадающие списки и др.)

- **backend/** — серверная часть на Go (Golang)
  - **internal/models/** — структуры данных для GORM (User, Canal, Client, Order и др.)
  - **internal/api/** — обработчики HTTP-запросов (handlers.go)
  - **cmd/** — точка входа приложения (main.go)
  - **routes/** — маршруты API (если используются)
  - **middleware/** — промежуточные обработчики (например, авторизация)
  - **db/** — инициализация и подключение к базе данных

## Пример структуры

```
ORDERS/
├── frontend/
│   └── lib/
│       ├── models/
│       ├── pages/
│       ├── services/
│       └── widgets/
├── backend/
│   ├── cmd/
│   ├── internal/
│   │   ├── models/
│   │   └── api/
│   ├── routes/
│   ├── middleware/
│   └── db/
```

## Основные файлы

- **frontend/lib/pages/0_0_login.dart** — страница входа пользователя
- **frontend/lib/pages/0_1_menu.dart** — главное меню
- **frontend/lib/pages/1_0_orders_create.dart** — создание заявки
- **frontend/lib/services/auth_service.dart** — работа с авторизацией и токеном
- **backend/internal/models/models.go** — структуры данных для GORM
- **backend/internal/api/handlers.go** — обработчики HTTP-запросов (API endpoints)
- **backend/cmd/main.go** — запуск сервера

## Как искать нужный код

- **API-запросы и авторизация** — ищи в `frontend/lib/services/`
- **UI и формы** — ищи в `frontend/lib/pages/` и `frontend/lib/widgets/`
- **Модели данных** — ищи в `frontend/lib/models/` и `backend/internal/models/`
- **Бизнес-логика и работа с БД** — ищи в `backend/internal/api/` и `backend/internal/models/`
- **Маршруты и запуск сервера** — ищи в `backend/cmd/` и `backend/routes/`

---

> **Совет:**  
> Используй поиск по проекту (`Ctrl+Shift+F` в VS Code) и переход к определению (`F12`), чтобы быстро находить нужные файлы и функции.

---

Если нужно добавить описание для конкретного модуля — просто дополни этот файл!

## Запуск проекта

- **backend (Go):**
  - Перейди в папку `backend`
  - Выполни команду:  
    ```
    go run ./cmd/main.go
    ```
- **frontend (Flutter):**
  - Перейди в папку `frontend`
  - Выполни команду:  
    ```
    flutter run
    ```

---
