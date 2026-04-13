Review the Go code in the provided file or selection. Act as a senior Go backend engineer performing a thorough code review.

## Review checklist

Analyze the code against these criteria and report findings:

### Correctness
- Logic errors, off-by-one, nil pointer dereference risks
- Error handling: all errors checked, wrapped with context, proper propagation
- Concurrency safety: race conditions, goroutine leaks, deadlocks

### Go Idioms
- Naming conventions (MixedCaps, short receivers, meaningful names)
- Early returns over deep nesting
- Accept interfaces, return structs
- Proper use of context.Context

### Security
- SQL injection, input validation, auth bypass risks
- Secrets in code, sensitive data in logs
- OWASP top 10 applicable issues

### Performance
- Unnecessary allocations, missing pre-allocation
- N+1 queries, missing indexes hints
- Unbounded operations (no pagination, no timeout)

### Design
- Single Responsibility Principle violations
- Interface bloat (too many methods)
- Proper layer separation (handler → service → repository)

## Output format

For each finding:
- **Severity**: 🔴 Critical / 🟡 Warning / 🔵 Suggestion
- **Location**: file:line
- **Issue**: what's wrong
- **Fix**: concrete code suggestion

End with a summary: total findings by severity, overall quality score (1-10), and top 3 action items.

$ARGUMENTS
