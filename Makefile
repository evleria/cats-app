swag:
	swag init --parseDependency --parseDepth=5
compose-build:
	docker-compose build
compose-up:
	docker-compose up -d
compose-down:
	docker-compose down
lint-backend:
	golangci-lint run ./backend/internal/... && golangci-lint run ./backend/main.go
import:
	goimports -local "github.com/evleria/mongo-crud" -w .

.PHONY: swag, compose-build, compose-up, compose-down, lint-backend, import