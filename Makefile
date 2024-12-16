export MIGRATION_FOLDER=$(INTERNAL_REPO_PATH)migrations

build:
	docker compose build

up-all:
	docker compose up -d

down:
	docker compose down

test:
	go test -coverprofile=c.out ./... -coverpkg="./..." && go tool cover -func c.out | grep total

migration-up-oms:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_OMS_SETUP)" up

migration-down-oms:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_OMS_SETUP)" down

BIN_DIR := ./bin

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

auth_service: $(BIN_DIR)
	go build -o $(BIN_DIR)/auth_service ./cmd/auth/main.go

user_service: $(BIN_DIR)
	go build -o $(BIN_DIR)/user_service ./cmd/user/main.go

event_service: $(BIN_DIR)
	go build -o $(BIN_DIR)/event_service ./cmd/event/main.go

image_service: $(BIN_DIR)
	go build -o $(BIN_DIR)/image_service ./cmd/image/main.go

notification_service: $(BIN_DIR)
	go build -o $(BIN_DIR)/notification_service ./cmd/notification/main.go

server_service: $(BIN_DIR)
	go build -o $(BIN_DIR)/server_service ./cmd/server/main.go

all: auth_service user_service event_service image_service notification_service server_service

clean:
	rm -rf $(BIN_DIR)
