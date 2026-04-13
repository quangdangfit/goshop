---
description: Build, test, mock generation, and Swagger doc commands
globs:
alwaysApply: true
---

# Commands

```bash
# Build
go build -o main cmd/api

# Run all tests with coverage
make unittest
# Equivalent: go test -timeout 9000s -a -v -coverprofile=coverage.out -coverpkg=./... ./...

# Run a single test package
go test ./internal/product/service/... -v -run TestProductServiceTestSuite

# Run a specific test case
go test ./internal/product/service/... -v -run TestProductServiceTestSuite/TestCreateSuccess

# Regenerate all mocks
make mock
# Equivalent: go generate ./...

# Regenerate Swagger docs
make doc
# Equivalent: swag fmt && swag init -g ./cmd/api/main.go
```
