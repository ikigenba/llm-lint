# Phase 27 — Swap the verbose line to id-first ordering

*Realizes design Decision 8 (verbose line shape) and Decision 1's restated
verbose ids.*

`report.VerboseSink` renders each progress line as
`<circle> [<rule-id>] <path>` — the bracketed rule id directly after the
circle, the path last — replacing the previous path-first ordering. The two
existing end-to-end tests update to the new shape.

**Done when:**

- R-H1KL-BCFY — `run` with `--verbose` writes one complete, non-interleaved
  `<circle> [<rule-id>] <path>` line per (rule, file) pair, circle 🟢/🔴 by
  verdict, bracketed rule id then file path, findings on `out` unchanged, no
  lines without the flag — its existing test updated to the new ordering.
- R-H2SH-P46N — `--verbose --format json` keeps the JSON findings on `out`
  and the `<circle> [<rule-id>] <path>` lines on `errOut`, unmixed — its
  existing test updated to the new ordering.
- The suite is green (per design CONVENTIONS.md).
