doc:
	swag fmt && swag init -g ./cmd/api/main.go

SOURCE_PKGS := $(shell go list ./... | grep -v '/mocks$$' | grep -v '/proto/gen/' | grep -v '/docs$$' | grep -v '/cmd/' | tr '\n' ',')

unittest:
	go test -timeout 9000s -v -coverprofile=coverage.out -coverpkg=$(SOURCE_PKGS) ./... 2>&1 | tee report.out

integration:
	go test -tags=integration -timeout 9000s -v -coverprofile=coverage.integration.out -coverpkg=$(SOURCE_PKGS) ./tests/integration/...

lint:
	golangci-lint run ./...

mock:
	go generate ./...

# Migrations
# Requires golang-migrate CLI: brew install golang-migrate
# (or: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
# Reads DSN from DATABASE_URI (falls back to the value in config.yaml is the caller's job).
DATABASE_URI ?= postgres://postgres:test@localhost:5432/goshop?sslmode=disable

migrate-up:
	migrate -path migrations -database "$(DATABASE_URI)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URI)" down 1

migrate-status:
	migrate -path migrations -database "$(DATABASE_URI)" version

# Scaffold a new migration: make migrate-new name=add_orders_index
migrate-new:
	@test -n "$(name)" || (echo "usage: make migrate-new name=<short_description>" && exit 1)
	migrate create -ext sql -dir migrations -seq $(name)
