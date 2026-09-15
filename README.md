# Barrage

A compact concurrent HTTP load tester implemented independently in Go.

## Features
- configurable request count and worker concurrency
- HTTP status/error accounting
- total throughput
- p50, p95 and p99 latency reporting
- standard-library-only implementation

```bash
go run . -url https://example.com -n 1000 -c 50
```

This repository is an independent implementation created for Adewale Babalola's systems-engineering portfolio. It does not copy the source of the similarly named upstream reference project.
