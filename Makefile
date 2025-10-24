POSTGRES_USER := postgres
POSTGRES_PASSWORD := postgres
DB_HOST := 127.0.0.1
DB_PORT := 5432
DB_NAME := master

TEST_DB_PORT := 5433
TEST_DB_NAME := postgres

DB_DSN="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"
TEST_DB_DSN="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(TEST_DB_PORT)/$(TEST_DB_NAME)?sslmode=disable"

API_SCHEMA="api/schema.yml"

.PHONY: api
api: ## Открыть схему контракта API в браузере
	@swgui ./api/schema.yml

.PHONY: codegen
codegen: ## Сгенерировать код по документации API
	@oapi-codegen -package dto -generate types -o internal/generated/dto.go $(API_SCHEMA)
	@oapi-codegen -package dto -generate gorilla -o internal/generated/sdk.go $(API_SCHEMA)

.PHONY: up
up: ## Поднять окружение
	@docker compose -f docker-compose.yml up -d

.PHONY: down
down: ## Поднять окружение
	@docker compose -f docker-compose.yml down

.PHONY: setup
setup: ## Установить все необходимые утилиты
	go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: migrations-up
migrations-up: ## Накатить миграции
	goose -dir db/migrations/master postgres $(DB_DSN) up

.PHONY: migrations-down
migrations-down: ## Откатить миграции
	goose -dir db/migrations/master postgres $(DB_DSN) down-to 0

.PHONY: migration-create
migration-create: ## Пример команды для создания миграции
	@echo "goose -dir db/migrations/<db> create <add_some_column> sql"

.PHONY: build
build: ## Сбилдить бинарь приложения
	go build -o ./bin/service ./cmd/service/main.go

.PHONY: testgen
testgen: ## Сгенерировать основу для тестов
	./bin/testgen -path "$(CURDIR)/$(PATH)"

.PHONY: test
test: ## Запустить юнит-тесты
	@go test -cover -shuffle=off ./...

ALLURE_RESULTS_DIR := $(PWD)

.PHONY: test-report
test-report: ## Запустить тесты с генерацией отчета Allure
	@rm -r $(ALLURE_RESULTS_DIR)/allure-results
	@mkdir -p $(ALLURE_RESULTS_DIR)
	@ALLURE_OUTPUT_PATH=$(ALLURE_RESULTS_DIR) go test ./...
	@allure generate ./allure-results -o ./allure-report --clean
	@allure open ./allure-report

.PHONY: test-random
test-random: ## Запустить юнит-тесты в рандомном порядке
	@go test -shuffle=on ./...

.PHONY: test-trace
test-trace: ## Запустить тесты с трейсингом по тредам и процессам
	@go test -trace=trace.out ./internal/usecase/cart
	@go tool trace trace.out

.PHONY: integration-test
integration-test: ## Запустить интеграционные тесты
	@docker compose -f tests/docker-compose.yml up -d
	@echo 'Starting environment...'
	@sleep 5s
	@echo 'Applying migrations...'
	@-goose -dir db/migrations/master postgres $(TEST_DB_DSN) up
	@rm -r $(ALLURE_RESULTS_DIR)/allure-results
	@mkdir -p $(ALLURE_RESULTS_DIR)
	@ALLURE_OUTPUT_PATH=$(ALLURE_RESULTS_DIR) go test -tags=integration ./tests/integration/...
	@allure generate ./allure-results -o ./allure-report --clean
	@allure open ./allure-report

.PHONY: e2e-test
e2e-test: ## Запустить e2e тесты
	@docker compose -f tests/docker-compose.e2e.yml up --build --abort-on-container-exit
	@docker compose -f tests/docker-compose.e2e.yml down
	@mkdir -p tests/allure-report
	@allure generate tests/allure-results -o tests/allure-report --clean
	@allure open tests/allure-report

.PHONY: integration-down
integration-down: ## Запустить интеграционные тесты
	@docker compose -f tests/docker-compose.yml down

.PHONY: cli
cli: ## Собрать бинарь для CLI
	go build -o shop cmd/cli/main.go

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL:=help
