export MIGRATION_FOLDER=$(INTERNAL_REPO_PATH)migrations
#
#build:
#	docker compose build
#
#up-all:
#	docker compose up -d
#
#down:
#	docker compose down
#
test:
	go test -coverprofile=c.out ./... -coverpkg="./..." && go tool cover -func c.out | grep total
#
#up-db:
#	docker compose up -d postgres
#
#stop-db:
#	docker compose stop postgres
#
#start-db:
#	docker compose start postgres
#
#down-db:
#	docker compose down postgres
#
migration-up-oms:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_OMS_SETUP)" up

migration-down-oms:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_OMS_SETUP)" down

# Имя Docker Compose файла
COMPOSE_FILE = docker-compose.yml

# Функция для запуска команды с передачей имени сервиса
define service_action
	@echo "Performing action '$(2)' for service: $(1)"
	docker compose -f $(COMPOSE_FILE) $(2) $(1)
endef

# Цель по умолчанию
.DEFAULT_GOAL := help

# Пересборка и запуск указанного сервиса с выводом логов
rebuild-service:
	@echo "Rebuilding and starting service: $(SERVICE)"
	docker compose -f $(COMPOSE_FILE) up --build $(SERVICE)

# Полная сборка и запуск всех сервисов с выводом логов
up:
	@docker compose -f $(COMPOSE_FILE) up --build

# Перезапуск указанного сервиса с выводом логов
restart-service:
	@echo "Restarting service: $(SERVICE)"
	docker compose -f $(COMPOSE_FILE) up $(SERVICE)

# Остановка указанного сервиса (принимает имя сервиса как аргумент)
stop-service:
	@$(call service_action,$(filter-out $@,$(MAKECMDGOALS)),stop)

# Полное удаление сервиса (принимает имя сервиса как аргумент)
rm-service:
	@$(call service_action,$(filter-out $@,$(MAKECMDGOALS)),rm -f)

# Просмотр логов для указанного сервиса (принимает имя сервиса как аргумент)
logs:
	@$(call service_action,$(filter-out $@,$(MAKECMDGOALS)),logs -f)

# Полная остановка всех сервисов
down:
	@docker compose -f $(COMPOSE_FILE) down

# Просмотр статуса всех сервисов
ps:
	@docker compose -f $(COMPOSE_FILE) ps

# Просмотр логов всех сервисов
logs-all:
	@docker compose -f $(COMPOSE_FILE) logs -f

# Удаление неиспользуемых данных
prune:
	@docker system prune -f --volumes

# Подсказка по использованию Makefile
help:
	@echo "Использование:"
	@echo "  make up                       - Запуск всех сервисов с пересборкой"
	@echo "  make down                     - Остановка всех сервисов"
	@echo "  make rebuild-service <service> - Пересборка и запуск одного сервиса с логами"
	@echo "  make restart-service <service> - Перезапуск одного сервиса с логами"
	@echo "  make stop-service <service>    - Остановка одного сервиса"
	@echo "  make rm-service <service>      - Удаление одного сервиса"
	@echo "  make logs <service>            - Просмотр логов одного сервиса"
	@echo "  make logs-all                 - Просмотр логов всех сервисов"
	@echo "  make ps                       - Статус всех сервисов"
	@echo "  make prune                    - Удаление неиспользуемых данных"

# Игнорирование аргументов как ошибок
%:
	@:

