Analyze the specified Go code or endpoint for performance issues and optimize.

## Analysis steps

1. **Identify the scope**: Which file, function, or endpoint to analyze?

2. **Check for common Go performance issues**:

   ### Memory
   - Unnecessary allocations in hot paths (loops, handlers)
   - Missing slice pre-allocation: `make([]T, 0, expectedLen)`
   - String concatenation in loops (use `strings.Builder`)
   - Large structs passed by value instead of pointer
   - Unreleased resources (unclosed readers, connections)

   ### Database
   - N+1 query problems (missing Preload/Joins)
   - Missing database indexes for WHERE/ORDER BY columns
   - Fetching all columns when only a few are needed (missing `Select()`)
   - Unbounded queries (missing LIMIT)
   - Unnecessary COUNT queries
   - Transaction held too long

   ### Concurrency
   - Goroutine leaks (no cancellation path)
   - Lock contention (mutex held during I/O)
   - Unbuffered channels causing blocking
   - Missing connection pool tuning

   ### Caching
   - Repeated identical DB queries that should be cached
   - Cache with no TTL or no invalidation
   - Over-caching (caching user-specific data globally)

   ### HTTP
   - Missing timeouts on outbound HTTP calls
   - Response body not closed
   - Large response bodies not streamed

3. **Suggest optimizations**: For each finding:
   - Explain the issue and its impact (latency, memory, CPU)
   - Provide concrete code fix
   - Estimate improvement (rough order of magnitude)

4. **Benchmark** (if applicable):
   - Write a benchmark: `func BenchmarkXxx(b *testing.B)`
   - Run: `go test -bench=. -benchmem ./path/...`
   - Compare before/after

$ARGUMENTS
