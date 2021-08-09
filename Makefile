swag:
	swag init --parseDependency --parseDepth=5
compose:
	docker-compose build && docker-compose up -d
compose-down:
	docker-compose down
lint:
	golangci-lint run ./internal/... && golangci-lint run ./main.go
import:
	goimports -local "github.com/evleria/mongo-crud" -w .
gen-mocks:
	mockery --all --recursive --inpackage --case underscore

.PHONY: swag, compose-build, compose-up, compose-down, lint, import, gen-mocks