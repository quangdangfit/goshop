Refactor the specified Go code following Go best practices and clean architecture principles.

## Process

1. **Analyze** the current code: read the file and all related files (callers, interfaces, tests).
2. **Identify** refactoring opportunities:
   - Functions > 50 lines → extract into smaller, well-named functions
   - Duplicated logic → extract shared helper or use composition
   - God structs → split by responsibility
   - Concrete dependencies → extract interfaces
   - Complex conditionals → use early returns, strategy pattern, or lookup tables
   - Magic numbers/strings → named constants
   - Poor naming → rename to match Go conventions
3. **Plan** the refactoring steps. Present the plan to me for approval before making changes.
4. **Execute** each step:
   - Make one logical change at a time
   - Update all callers and imports
   - Regenerate mocks if interfaces changed (`mockery`)
   - Run tests after each step to verify nothing breaks
5. **Verify**: run `go vet ./...` and `golangci-lint run` on affected packages.

## Constraints

- Do NOT change public API signatures unless explicitly approved
- Do NOT mix refactoring with feature changes
- Preserve all existing test coverage
- Keep backward compatibility with proto/gRPC contracts

$ARGUMENTS
