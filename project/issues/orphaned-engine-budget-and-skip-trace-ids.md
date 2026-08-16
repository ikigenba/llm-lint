Orphaned engine verification ids observed while resolving Phase 20's engine surface.

`internal/engine/context_budget_skip_test.go` is tagged `R-GWLC-EV7X`. `internal/engine/oversize_skip_trace_test.go` is tagged `R-H40E-2VXC`. Neither id appears in `project/design/INDEX.md` or Decision 5's Verification list. Decision 5 rejected pre-sizing files against a context window; the oversize test also constructs the buffering `engine.Trace` / `Engine.Trace` that Phase 20 is deleting.

Out of scope for this gather turn, which only briefs Phase 20's live `--verbose` slice (`R-H1KL-BCFY`, `R-H2SH-P46N`). Noting the orphaned ids so a later turn can delete or re-home the tests.
