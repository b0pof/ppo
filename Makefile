POSTGRES_USER := postgres
POSTGRES_PASSWORD := postgres
DB_HOST := 127.0.0.1

DB_PORT := 5432
DB_NAME := master

TEST_DB_PORT := 5433
TEST_DB_NAME := postgres

APP_PORT := 8080

DB_DSN="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"
TEST_DB_DSN="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(TEST_DB_PORT)/$(TEST_DB_NAME)?sslmode=disable"

API_SCHEMA="api/schema.yml"

.PHONY: api
api: ## Открыть схему контракта API в браузере
	@swgui $(API_SCHEMA)

.PHONY: httpcli
httpcli: ## Сбилдить бинарь для httpcli
	@go build -o httpcli cmd/httpcli/main.go

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

.PHONY: unit-test
unit-test: ## Запустить тесты с генерацией отчета Allure
	@go test -json ./internal/... | golurectl -l -e -o report/unit-allure-results --allure-suite Unit --allure-tags UNIT

.PHONY: report-merge
report-merge: ## Подготовить финальный отчет
	@rm -rf report/allure-results
	@mkdir -p report/allure-results
	@cp -r report/unit-allure-results/* report/allure-results/ 2>/dev/null || echo "No unit results"
	@cp -r report/integration-allure-results/* report/allure-results/ 2>/dev/null || echo "No integration results"
	@cp -r report/e2e-allure-results/* report/allure-results/ 2>/dev/null || echo "No e2e results"
	@echo "Total files: $(shell find report/allure-results -type f 2>/dev/null | wc -l)"

.PHONY: report-open
report-open: ## Открыть отчет
	@make report-merge
	@allure generate ./report/allure-results -o ./report/report --clean
	@allure open ./report/report

.PHONY: test-random
test-random: ## Запустить юнит-тесты в рандомном порядке
	@go test -shuffle=on ./...

.PHONY: test-trace
test-trace: ## Запустить тесты с трейсингом по тредам и процессам
	@go test -trace=trace.out ./internal/usecase/cart
	@go tool trace trace.out

.PHONY: integration-test
integration-test: ## Запустить интеграционные тесты
	@make deps
	@docker compose -f tests/docker-compose.yml up -d
	@echo 'Starting environment...'
	@sleep 3s
	@echo 'Applying migrations...'
	@-goose -dir db/migrations/master postgres $(TEST_DB_DSN) up
	@go test -p 1 -json -tags=integration ./tests/integration/... | golurectl -l -e -o report/integration-allure-results --allure-suite Integration --allure-tags INTEGRATION

.PHONY: e2e-test
e2e-test: ## Запустить e2e тесты
	@make deps
	@lsof -ti:$(APP_PORT) | xargs -r kill -9
	@docker compose -f tests/docker-compose.e2e.yml up -d
	@sleep 3
	@-goose -dir db/migrations/master postgres $(TEST_DB_DSN) up
	@sleep 3
	@POSTGRES_PORT=5433 POSTGRES_DATABASE=postgres go run cmd/service/main.go & \
		sleep 3; \
		go test -json -p 1 -tags=e2e ./tests/e2e/... | golurectl -l -e -o report/e2e-allure-results --allure-suite E2E --allure-tags E2E; \
		TEST_EXIT_CODE=$$?; \
		lsof -ti:$(APP_PORT) | xargs -r kill -9
		exit $$TEST_EXIT_CODE

.PHONY: e2e-down
e2e-down: ## Удалить окружение е2е тестов
	@docker compose -f tests/docker-compose.e2e.yml down

.PHONY: integration-down
integration-down: ## Удалить окружение интеграционных тестов
	@docker compose -f tests/docker-compose.yml down

.PHONY: deps
deps: ## Установить зависимости
	@go install github.com/robotomize/go-allure/cmd/golurectl@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: cli
cli: ## Собрать бинарь для CLI
	go build -o shop cmd/cli/main.go

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL:=help
