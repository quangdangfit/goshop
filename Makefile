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
