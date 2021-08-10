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
protoc:
	 rm -rf ./internal/pb && protoc --proto_path=proto proto/*.proto --go_out=./internal --go-grpc_out=./internal
grpcui:
	grpcui -plaintext localhost:$(PORT)

.PHONY: swag, compose-build, compose-up, compose-down, lint, import, gen-mocks, protoc, grpcui