# Phase 6 — Verdict cache decorator

*Realizes design Decision 6 (Verdict cache). Depends on Phase 5.*

End state: `internal/cache` per D06 — content-hash `Key`, fan-out JSON store,
`CachingClient` decorator with `NoRead`, hit counting, miss-and-continue
degradation on cache I/O errors, and mtime-based 30-day pruning with an
injected clock.

**Done when:** R-H1GX-XY6P, R-H2OU-BPXE, R-H3WQ-PHO3, R-H54N-39ES,
R-H6CJ-H15H are each discharged by a tagged test in `internal/cache`, and the
suite is green.
