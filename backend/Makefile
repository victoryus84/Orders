.PHONY: up down restart ps build clean help

# Включает переменные окружения из .env (должен быть в формате Makefile: VAR=value без пробелов вокруг =)
include .env
# Export all variables from .env for docker compose
export

# Запуск всего стека 
up:
	docker compose -f compose-prod.yml up -d

# Остановка всего стека
down:
	docker compose -f compose-prod.yml down --remove-orphans

# Полная перезагрузка
restart: 
	docker compose -f compose-prod.yml down --remove-orphans
	docker compose -f compose-prod.yml up -d

# Очистка
# ВНИМАНИЕ: команда удаляет named volumes, что может привести к потере данных!
clean:
	docker compose -f compose-prod.yml down --remove-orphans -v
	docker system prune -f

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  make up          - Запуск всего стека"
	@echo "  make down        - Остановка стека"
	@echo "  make restart     - Перезагрузка стека"
	@echo "  make clean       - Полная очистка (ВНИМАНИЕ: удаляет named volumes и может привести к потере данных!)"
	@echo "  make help        - Эта справка"