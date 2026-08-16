# Phase 19 — Reshape the Client seam: real usage, drop oversize-skip and OpFailure

*Realizes design Decision 5 (inference engine) and 8 (reporting, stats/verbose),
with the docs-level error-surface change in Decision 1 (CLI surface). No
dependency on a pending phase — the queue is otherwise empty.*

This phase reshapes the inference seam so `--stats` reports real token and cost
figures, and removes two mechanisms the redesign retired: the context-window
oversize-skip path and the `engine.OpFailure` typed error. It is a single
coherent change: the `Client.Judge` signature moves, which ripples through
`engine`, `cache`, `report`, `cmd/llm-lint`, and every fake Client in the
tests, so splitting it would leave the tree non-compiling between phases.

Observable end state:

- **Usage on the seam.** `engine.Usage{Input, Output int64; CostUSD float64}`
  exists; `engine.Client.Judge` returns `([]Finding, Usage, error)`.
  `AgentkitClient.Judge` fills `Usage` from the drained stream's `Usage()` and
  `Cost().USD()` (Input from the input buckets, Output including reasoning).
  `cache.CachingClient.Judge` returns the zero `Usage` on a cache hit and
  forwards `Next`'s `Usage` on a miss. `Engine.Run` sums each pair's `Usage`
  into `Stats.InputTokens`, `Stats.OutputTokens`, and `Stats.CostUSD`; the
  `--stats` line renders the run's real totals instead of zero placeholders.
- **Oversize-skip removed.** `AgentkitClient.Context`, `ContextWindow()`, the
  engine's pre-call size check and its `skipped` trace path are gone;
  `TraceEntry` carries only `pass`/`fail` outcomes. The tests that exercised
  the skip (`internal/engine/context_budget_skip_test.go`,
  `internal/engine/oversize_skip_trace_test.go`) are deleted along with their
  retired ids R-GWLC-EV7X and R-H40E-2VXC. An oversized file now surfaces as
  an ordinary operational failure (exit 3).
- **OpFailure removed.** `Engine.Run` returns a plain `%w`-wrapped `error`;
  the `OpFailure` type is gone. `run` still maps any engine error to exit 3
  (R-FYG5-PAGH unchanged).

**Done when:**

- R-ETNZ-VY1Q — a fake Client returning a known non-zero `Usage` per call over
  K judged pairs yields `Stats` whose `InputTokens`, `OutputTokens`, and
  `CostUSD` each equal K× the per-call value (non-zero) — proving per-pair
  usage is summed, not dropped. Realized by a genuine `internal/engine` test
  tagged with the id.
- R-HG3Q-J731 — reworked so its `--stats` test uses a fake Client reporting
  non-zero per-call usage and asserts the line's token and cost figures are the
  run's non-zero totals, no longer the zero placeholders.
- R-H7O3-875F — reworked so `Verbose` no longer renders (or asserts) a
  `skipped` line; only the judged-pair `<rule-id> <hit|miss> <pass|fail>`
  form remains.
- R-GZ15-6EPB — the live Google smoke additionally runs with `--stats` and
  asserts the emitted stderr line's input-token figure is non-zero after the
  real call (still skips without `GEMINI_API_KEY`).
- The suite is green: `go build ./...` and `go vet ./...` clean, `go test ./...`
  exits 0 (with no network and no credentials — the live/subscription smokes
  skip), and `gofmt -l .` prints nothing (excluding `project/`). No test in the
  tree references R-GWLC-EV7X or R-H40E-2VXC.
