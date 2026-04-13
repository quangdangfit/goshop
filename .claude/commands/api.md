Scaffold a new API endpoint for this Go project following the established architecture.

## Steps

1. **Gather requirements**: Parse the endpoint description from the arguments. Determine:
   - Domain: which domain does this belong to? (user, product, order, cart)
   - Resource name and HTTP method(s)
   - Request/response DTOs
   - Auth required? Which roles?
   - gRPC equivalent needed?

2. **Generate in order**:

   a. **Model** (`internal/{domain}/model/`) — if a new entity is needed:
      - GORM struct with proper tags, `BeforeCreate` UUID hook
      - Add to `AutoMigrate` in `cmd/api/main.go`

   b. **DTO** (`internal/{domain}/dto/`) — request and response structs:
      - Validation tags on request fields
      - Keep response flat, no nested models

   c. **Repository interface + implementation** (`internal/{domain}/repository/`):
      - Add methods to existing interface, or create new one
      - Use `dbs.Database` with functional options
      - Update `//go:generate mockery` comment if new interface

   d. **Service interface + implementation** (`internal/{domain}/service/`):
      - Business logic, input validation, error mapping
      - Depend on repository interface, not concrete

   e. **HTTP handler** (`internal/{domain}/port/http/`):
      - Gin handler with Swagger annotations
      - Proper error responses with correct HTTP status codes
      - Register routes in the routes function

   f. **gRPC handler** (`internal/{domain}/port/grpc/`) — if needed:
      - Proto message definitions
      - Handler implementation
      - Register in server

3. **Generate mocks**: Run `mockery` for new/updated interfaces.
4. **Generate tests**: Create handler and service test suites.
5. **Update Swagger**: Run `make doc`.
6. **Verify**: Run `go build ./...` and `go test ./...`.

$ARGUMENTS
