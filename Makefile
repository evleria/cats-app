swag:
	swag init --parseDependency --parseDepth=5
compose:
	docker-compose build && docker-compose up -d
compose-down:
	docker-compose down
lint:
	golangci-lint run ./internal/... && golangci-lint run ./main.go
import:
	goimports -local "github.com/evleria/cats-app" -w .
gen-mocks:
	mockery --all --recursive --inpackage --case underscore
protoc:
	 rm -rf ./protocol/pb && protoc --proto_path=protocol/proto protocol/proto/*.proto --go_out=./protocol --go-grpc_out=./protocol
grpcui:
	grpcui -plaintext localhost:$(PORT)

.PHONY: swag, compose-build, compose-up, compose-down, lint, import, gen-mocks, protoc, grpcui