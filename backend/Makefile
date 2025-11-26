.PHONY: ssl-init ssl-renew ssl-status ssl-test ssl-logs up down restart ps build clean help

# Включает переменные окружения из .env (должен быть в формате Makefile: VAR=value без пробелов вокруг =)
include .env
# Export all variables from .env for docker compose
export

# Первоначальная настройка SSL
ssl-init:
	docker compose up certbot-init

# Ручное обновление сертификатов
ssl-renew:
	docker compose exec certbot-renew /usr/local/bin/certbot-renew.sh

# Проверка статуса сертификатов
ssl-status:
	docker compose exec certbot-renew /usr/local/bin/certbot-status.sh

# Просмотр логов обновления
ssl-logs:
	docker compose logs certbot
# Тестирование SSL
ssl-test:
	@if [ -z "$(DOMAIN)" ]; then \
		echo "Ошибка: переменная DOMAIN не установлена. Проверьте .env или задайте DOMAIN вручную."; \
		exit 1; \
	fi
	docker compose run --rm curlimages/curl curl -I "https://${DOMAIN}/health"
	docker compose run --rm curlimages/curl curl -I "https://${DOMAIN}/health"

# Сборка образов
build:
	docker compose build certbot-init certbot-renew

# Запуск всего стека (без certbot-init)
up:
	docker compose up -d postgres api nginx certbot

# Остановка всего стека
down:
	docker compose down --remove-orphans

# Полная перезагрузка
restart: down up

# Очистка
# ВНИМАНИЕ: команда удаляет named volumes, что может привести к потере данных!
clean:
	docker compose down --remove-orphans -v
	docker system prune -f
# Очистка
clean:
	docker compose down --remove-orphans -v
	docker system prune -f

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  make ssl-init    - Первоначальная настройка SSL"
	@echo "  make ssl-renew   - Ручное обновление сертификатов"
	@echo "  make ssl-status  - Проверка статуса сертификатов"
	@echo "  make ssl-test    - Тестирование SSL"
	@echo "  make clean       - Полная очистка (ВНИМАНИЕ: удаляет named volumes и может привести к потере данных!)"
	@echo "  make down        - Остановка стека"
	@echo "  make restart     - Перезагрузка стека"
	@echo "  make clean       - Полная очистка (ВНИМАНИЕ: 'docker system prune -f' удаляет неиспользуемые данные Docker на всей системе, что может повлиять на другие проекты)"
	@echo "  make help        - Эта справка"ка"
	@echo "  make help        - Эта справка"