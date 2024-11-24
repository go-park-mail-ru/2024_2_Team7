export MIGRATION_FOLDER=$(INTERNAL_REPO_PATH)migrations

build:
	docker compose build

up-all:
	docker compose up -d 

down:
	docker compose down

test:
	go test -coverprofile=c.out ./... -coverpkg="./..." && go tool cover -func c.out | grep total

up-db:
	docker compose up -d postgres

stop-db:
	docker compose stop postgres

start-db:
	docker compose start postgres

down-db:
	docker compose down postgres

migration-up-oms:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_OMS_SETUP)" up

migration-down-oms:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_OMS_SETUP)" down