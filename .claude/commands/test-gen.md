Generate comprehensive Go tests for the specified file, function, or package.

## Instructions

1. Read the target code thoroughly. Understand all inputs, outputs, error paths, and edge cases.
2. Follow the project's testing conventions:
   - Use `testify/suite` with `SetupTest` for service and handler tests
   - Use table-driven tests with `t.Run` for utility functions
   - Mock dependencies using mockery-generated mocks from `{package}/mocks/`
   - Handler tests use `gin.CreateTestContext` + `httptest.ResponseRecorder`
3. Generate tests covering:
   - **Happy path**: normal successful operation
   - **Validation errors**: invalid input, missing required fields
   - **Not found**: entity doesn't exist
   - **Authorization**: wrong user/role accessing resource
   - **Edge cases**: empty lists, zero values, nil pointers, duplicate entries
   - **Error propagation**: DB errors, external service failures
4. Use realistic test data, not "test123" placeholders.
5. Assert on:
   - Return values and error types
   - HTTP status codes and response body for handlers
   - Mock expectations (`.AssertExpectations(t)`)
6. Run the generated tests to verify they compile and pass.

## Naming convention
- Suite: `Test{Type}TestSuite` (e.g. `TestProductServiceTestSuite`)
- Methods: `Test{Action}{Scenario}` (e.g. `TestCreateSuccess`, `TestCreateValidationError`)

$ARGUMENTS
