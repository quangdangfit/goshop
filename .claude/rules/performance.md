---
description: Performance optimization and concurrency patterns for Go backend services
globs:
  - "**/*.go"
alwaysApply: false
---

# Performance & Concurrency

## Concurrency patterns

- Use `context.Context` for cancellation and timeouts on all I/O operations (DB queries, HTTP calls, Redis).
- Goroutines: always ensure they can exit. Pass a context or use `done` channels. Never fire-and-forget goroutines that can leak.
- Use `sync.WaitGroup` for fan-out/fan-in. Use `errgroup.Group` when goroutines can return errors.
- Protect shared state with `sync.Mutex` or `sync.RWMutex`. Prefer channels for communication, mutexes for state protection.
- Use `sync.Once` for expensive one-time initialization (e.g. DB connections, config loading).

## Database performance

- Always use database indexes for columns in WHERE, JOIN, and ORDER BY clauses.
- Use `Preload` selectively — only preload associations that are actually needed.
- Paginate all list endpoints. Never return unbounded result sets.
- Use `Select()` to fetch only needed columns for read-heavy endpoints.
- Batch operations: use `CreateInBatches` for bulk inserts, transactions for multi-step writes.
- Connection pooling: configure `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime` on the DB pool.

## Caching

- Cache frequently read, rarely written data (product listings, categories).
- Use cache-aside pattern: check cache → miss → query DB → populate cache.
- Set appropriate TTLs. Invalidate on writes using pattern-based deletion.
- Never cache user-specific or sensitive data without proper scoping.

## Memory & allocations

- Pre-allocate slices when the size is known: `make([]T, 0, expectedSize)`.
- Use `strings.Builder` for string concatenation in loops, not `+`.
- Use `sync.Pool` for frequently allocated/deallocated objects in hot paths.
- Avoid unnecessary pointer indirection — small structs are cheaper to copy.

## HTTP performance

- Use connection pooling for outbound HTTP clients. Reuse `http.Client` instances.
- Set timeouts on all HTTP clients: `Timeout`, `DialContext` deadline.
- Use streaming for large payloads instead of buffering everything in memory.
- Enable gzip compression for responses > 1KB.

## Profiling

- Use `pprof` for CPU and memory profiling in development.
- Benchmark critical paths with `go test -bench=. -benchmem`.
- Monitor goroutine count in production to detect leaks.
