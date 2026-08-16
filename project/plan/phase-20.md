# Phase 20 — Stream the --verbose audit trace live

*Realizes design Decision 8 (VerboseSink), 5 (TraceSink), 6 (caching client
write-through), and 1 (run wiring). Depends on Phase — (extends the built tree).*

Reshape `--verbose` from a buffered, sorted end-of-run dump into a live stream:
each judged (rule, file) pair's audit line is written to stderr the moment its
verdict is known, in completion order, with no post-run render.

Observable end state:

- `internal/engine`: the buffering `Trace` struct (`Add`/`Entries`) is gone,
  replaced by the `TraceSink` interface (`Add(TraceEntry)`); `TraceEntry`
  stays. The dead `Trace` field is removed from `Engine`.
- `internal/report`: `Verbose(w, cwd, root, entries)` (sort-then-render) is
  gone, replaced by `VerboseSink` (`NewVerboseSink(w, cwd, root)` + `Add`), a
  live `engine.TraceSink` that renders each entry to one
  `<path>: <rule-id> <hit|miss> <pass|fail>` line (path relative to cwd, no
  `llm-lint:` prefix) and holds a mutex across each write so concurrent lines
  never interleave; a write error on the audit stream is swallowed.
- `internal/cache`: `CachingClient.Trace` is now an `engine.TraceSink`; each
  judged pair calls `Trace.Add` the moment its verdict is known (unchanged
  Cached/Outcome semantics), streaming rather than accumulating.
- `cmd/llm-lint`: on `--verbose`, `run` builds `report.NewVerboseSink(errOut,
  cwd, cfg.Root)` and injects it as the caching client's `Trace`; it no longer
  constructs an `engine.Trace`, passes a trace to `engine.Engine`, or calls
  `report.Verbose` after the run. On exit 3 the already-streamed lines stand.
- The `report` sort/format unit test tagged `R-H7O3-875F`
  (`internal/report/verbose_sort_format_test.go`) is deleted with the behavior
  it proved; the D6 cache trace tests migrate from `*engine.Trace` to a fake
  `engine.TraceSink`; the D1 tests tagged `R-H1KL-BCFY` and `R-H2SH-P46N` are
  rewritten to assert the streaming contract below.

**Done when:** `go build ./...`, `go vet ./...`, `go test ./...` exit 0 and
`gofmt -l .` (excluding `project/`) is empty, with these ids covered by genuine
tests:

- `R-H1KL-BCFY` — `run` with `--verbose` (fake engine/cache over ≥2 rules ×
  ≥2 files, concurrency > 1) writes to `errOut` one line per (rule, file) pair,
  each a complete well-formed `<path>: <rule-id> <hit|miss> <pass|fail>` line
  (no two writes interleave), carrying the pair's cache hit/miss and pass/fail;
  the lines match as an order-independent set (completion order, not sorted);
  `out` carries only the findings it would have without the flag; and the same
  run without `--verbose` writes none of those lines to `errOut`.
- `R-H2SH-P46N` — `run` with `--verbose --format json` writes the findings as
  the JSON array to `out` and the per-pair audit lines as plain text to
  `errOut`; the two streams do not intermix.
- No test tagged `R-H7O3-875F` remains anywhere
  (`grep -rn 'R-H7O3-875F' --include='*.go' .` is empty).
