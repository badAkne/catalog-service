OUTPUT:=./bin/app
GO_LINT_VERSION=1.64.8

GO_FILE:=./main.go

.PHONY: up
up: ## Поднимает (запускает) окружение для работы приложения
	docker compose up -d

.PHONY: down
down: ## Отключает окружение для работы приложения
	docker compose down --remove-orphans

.PHONY: lint
lint: ## Запуск линтера
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v${GO_LINT_VERSION} run

.PHONY: lint-fix
lint-fix: ## Запуск линтера с фиксом
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v${GO_LINT_VERSION} run --fix

.PHONY: build
build: ## Сборка приложения
	go build -o ${OUTPUT} ${GO_FILE}

.PHONY: test
test: ## Запуск тестов
	go test -count=1 -v ./...

.PHONY: cover
cover: ## Просмотр покрытия
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: generate-mocks
generate-mocks: ##Генерация моков
	chmod +x scripts/generate_mocks.sh
	export PATH=$$PATH:$$HOME/go/bin && ./scripts/generate_mocks.sh