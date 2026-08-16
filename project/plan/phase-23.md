# Phase 23 — Delete the orphaned engine skip-path guard tests

*Realizes design Decision 5 (inference engine) — structural cleanup, no ids.*

The engine package carries two relic tests guarding the absence of the
pre-sizing/skip feature Decision 5 rejected:
`internal/engine/context_budget_skip_test.go` (tagged R-GWLC-EV7X) and
`internal/engine/oversize_skip_trace_test.go` (tagged R-H40E-2VXC). Neither id
is minted by any Decision; their real content (one Judge call per pair) is
already proven by R-GRPQ-VS95. Delete both files. The end state is an engine
package whose tests carry only minted ids.

**Done when:**

- `grep -rnE 'R-GWLC-EV7X|R-H40E-2VXC' --exclude-dir=project .` prints
  nothing.
- The suite is green (per design CONVENTIONS.md).
