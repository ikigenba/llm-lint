# Phase 22 — Bracket the rule id on the --verbose line

*Realizes design Decision 8 (VerboseSink render) and 1 (run wiring, verbose
verification). Depends on Phase — (extends the built tree).*

Wrap the rule id in square brackets on the `--verbose` progress line, changing
the rendered format from `<circle> <path> <rule-id>` to
`<circle> <path> [<rule-id>]` so the rule stays legible after a long path.
Rendering only — the pass/fail circle, path, streaming, completion-order,
non-interleaving, and stderr-vs-stdout contracts are unchanged.

Observable end state:

- `internal/report`: `VerboseSink.Add` renders each entry to one
  `<circle> <path> [<rule-id>]` line — 🟢/🔴 from `TraceEntry.Outcome`, then the
  path relative to cwd, then the rule id in square brackets, single-space-
  separated with no colon and no `llm-lint:` prefix. The mutex-across-write
  non-interleaving and swallowed write error are unchanged; `TraceEntry.Cached`
  is still not read.
- `cmd/llm-lint`: no wiring change.
- The D1 tests tagged `R-H1KL-BCFY` and `R-H2SH-P46N` are updated to assert the
  bracketed format below; any `internal/report` `VerboseSink` render unit test
  is updated to match.

**Done when:** `go build ./...`, `go vet ./...`, `go test ./...` exit 0 and
`gofmt -l .` (excluding `project/`) is empty, with these ids covered by genuine
tests:

- `R-H1KL-BCFY` — `run` with `--verbose` (fake engine/cache over ≥2 rules ×
  ≥2 files, concurrency > 1) writes to `errOut` one line per (rule, file) pair,
  each a complete well-formed `<circle> <path> [<rule-id>]` line whose leading
  circle is 🟢 for a passing pair and 🔴 for a failing one, followed by the file
  path then the rule id wrapped in square brackets, single-space-separated with
  no colon and no cache hit/miss token; the lines are order-independent
  (completion order, not sorted); `out` carries only the findings it would have
  without the flag; and the same run without `--verbose` writes none of those
  lines to `errOut`.
- `R-H2SH-P46N` — `run` with `--verbose --format json` writes the findings as
  the JSON array to `out` (unchanged by the flag) and the per-pair
  `<circle> <path> [<rule-id>]` progress lines to `errOut`; the two streams do
  not intermix.
