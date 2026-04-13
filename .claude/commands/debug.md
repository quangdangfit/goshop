Diagnose and fix the described bug or issue in this Go codebase.

## Debugging process

1. **Reproduce**: Understand the issue from the description. Identify:
   - What is the expected behavior?
   - What is the actual behavior?
   - Which endpoint/function/flow is affected?

2. **Trace the flow**: Follow the code path from entry point to failure:
   - HTTP handler → service → repository → database
   - Check request binding, validation, business logic, DB queries
   - Look for nil pointer dereferences, incorrect type assertions, wrong error handling

3. **Narrow down**: Use these techniques:
   - Read error messages carefully — they usually point to the exact issue
   - Check recent git changes: `git log --oneline -20` and `git diff HEAD~5`
   - Search for related patterns: similar code that works vs. the broken code
   - Run specific test: `go test ./path/to/package/... -v -run TestName`
   - Check for common Go pitfalls:
     - Loop variable capture in goroutines
     - Nil map/slice operations
     - Interface nil vs typed nil
     - Context cancellation not checked
     - Race conditions (run with `-race` flag)

4. **Fix**: Apply the minimal fix that resolves the issue:
   - Fix the root cause, not the symptom
   - Don't refactor surrounding code — fix only what's broken
   - Add a test that reproduces the bug BEFORE fixing

5. **Verify**:
   - Run the failing test to confirm the fix
   - Run the full test suite to check for regressions
   - Run `go vet ./...` on affected packages

$ARGUMENTS
