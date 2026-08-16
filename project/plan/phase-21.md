# Phase 21 — Reshape the --verbose line to a pass/fail circle

*Realizes design Decision 8 (VerboseSink render) and 1 (run wiring, verbose
verification). Depends on Phase — (extends the built tree).*

Redefine the `--verbose` progress line from
`<path>: <rule-id> <hit|miss> <pass|fail>` to `<circle> <path> <rule-id>`: a
leading pass/fail emoji circle, then the file path, then the rule id, single
spaces, no colon, no cache hit/miss token. Rendering only — the streaming,
completion-order, non-interleaving, and stderr-vs-stdout contracts are
unchanged.

Observable end state:

- `internal/report`: `VerboseSink.Add` renders each entry to one
  `<circle> <path> <rule-id>` line — 🟢 when `TraceEntry.Outcome == "pass"`, 🔴
  when `== "fail"`, then the path relative to cwd, then the rule id,
  single-space-separated with no colon and no `llm-lint:` prefix. It no longer
  reads `TraceEntry.Cached` and emits no `hit`/`miss` token. The mutex-across-
  write non-interleaving and swallowed write error are unchanged. The emoji
  carry their own color, so there is no TTY detection and no `NO_COLOR` branch.
- `internal/engine`: `TraceEntry` is untouched — `Cached` remains on the seam
  (the caching client still records it); only the sink stops rendering it.
- `cmd/llm-lint`: no wiring change — `run` still builds
  `report.NewVerboseSink(errOut, cwd, cfg.Root)` on `--verbose` and injects it
  as the caching client's `TraceSink`.
- The D1 tests tagged `R-H1KL-BCFY` and `R-H2SH-P46N` are rewritten to assert
  the `<circle> <path> <rule-id>` format below (leading 🟢/🔴, no hit/miss
  token); any `internal/report` `VerboseSink` render unit test is updated to
  the new format.

**Done when:** `go build ./...`, `go vet ./...`, `go test ./...` exit 0 and
`gofmt -l .` (excluding `project/`) is empty, with these ids covered by genuine
tests:

- `R-H1KL-BCFY` — `run` with `--verbose` (fake engine/cache over ≥2 rules ×
  ≥2 files, concurrency > 1) writes to `errOut` one line per (rule, file) pair,
  each a complete well-formed `<circle> <path> <rule-id>` line whose leading
  circle is 🟢 for a passing pair and 🔴 for a failing one, followed by the file
  path then the rule id, single-space-separated with no colon and no cache
  hit/miss token; the lines are order-independent (completion order, not
  sorted); `out` carries only the findings it would have without the flag; and
  the same run without `--verbose` writes none of those lines to `errOut`.
- `R-H2SH-P46N` — `run` with `--verbose --format json` writes the findings as
  the JSON array to `out` (unchanged by the flag) and the per-pair
  `<circle> <path> <rule-id>` progress lines to `errOut`; the two streams do
  not intermix.
