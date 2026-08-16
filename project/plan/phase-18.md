# Phase 18 — `--verbose` per-pair audit trace

*Realizes design Decision 1 (CLI/run wiring), 5 (engine `Trace`/`TraceEntry`
and skip recording), 6 (caching-client judged-pair recording), and 8
(`report.Verbose` renderer).*

Adds the `--verbose` flag: on request, a run prints on stderr one line per
(rule, file) pair evaluated — cache hit/miss and pass/fail — leaving the
findings on stdout untouched. The pieces interlock through one shared trace
type, so they land together in a single phase rather than in half-built
states.

What gets built (observable end state):

- `internal/engine`: a concurrency-safe `Trace` type with `Add`/`Entries` and
  the `TraceEntry{File, Rule string; Cached bool; Outcome string}` record
  (D05). `Engine` gains a `Trace *Trace` field; when set, each oversize skip
  adds a `TraceEntry{Outcome: "skipped"}` and no `Judge` call is made for it.
- `internal/cache`: `CachingClient` gains a `Trace *engine.Trace` field; when
  set, each judged pair adds an entry carrying `Cached` (hit vs miss, always
  miss under `NoRead`) and `Outcome` `"pass"`/`"fail"` from the **raw**
  pre-suppression verdict.
- `internal/report`: `Verbose(w, cwd, root, entries)` sorts entries by
  (file, rule) and writes, per line, the cwd-relative path plus
  `<rule-id> <hit|miss> <pass|fail>` for judged pairs and `<rule-id> skipped`
  for skips, with no `llm-lint:` prefix; independent of `--format`.
- `cmd/llm-lint`: the `--verbose` flag; `run` creates one shared
  `*engine.Trace`, injects it into the engine and the caching client when the
  flag is set (nil otherwise), and after a completed run (exit 0 or 1) renders
  it with `report.Verbose` to `errOut`. On an operational failure (exit 3) no
  trace is rendered.

**Done when:** the suite is green (`go build ./...`, `go vet ./...`,
`go test ./...` exit 0, `gofmt -l .` excluding `project/` empty) and each id
below is discharged by a genuine test tagged with it verbatim:

- R-H1KL-BCFY — `run --verbose` over a fake engine/cache with ≥2 rules × ≥2
  files writes one per-pair line (hit/miss + pass/fail) per pair to `errOut`
  and only the findings to `out`; without the flag those lines are absent.
- R-H2SH-P46N — `run --verbose --format json` writes the JSON findings array
  to `out` and the plain-text per-pair lines to `errOut`, without intermixing.
- R-H40E-2VXC — with a `Trace` set, an oversize-skipped pair adds exactly one
  `TraceEntry{Outcome: "skipped"}` and makes no `Judge` call; with `Trace` nil
  it adds nothing.
- R-H58A-GNO1 — a warm-cache hit records `Cached=true`; a miss that calls
  `Next` records `Cached=false`.
- R-H6G6-UFEQ — a ≥1-finding verdict records `Outcome="fail"`, a clean one
  `"pass"`, and a pair whose sole finding is later suppressed by D07 still
  records `fail` (raw pre-suppression verdict).
- R-H7O3-875F — `Verbose` renders entries supplied out of order sorted by
  (file, rule), with cwd-relative paths, judged pairs as
  `<rule-id> <hit|miss> <pass|fail>` and skips as `<rule-id> skipped`.
