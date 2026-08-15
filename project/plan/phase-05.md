# Phase 5 — Inference engine against fakes

*Realizes design Decision 5 (Inference engine), slice R-GRPQ-VS95,
R-GSXN-9JZU, R-GU5J-NBQJ, R-GVDG-13H8, R-GWLC-EV7X, R-GXT8-SMYM (the live
ids belong to Phase 9). Depends on Phases 2 and 3.*

End state: `internal/engine` per D05 — the `Client` seam,
`AgentkitClient.Judge` with the `report_violations` typed tool, the
one-retry-then-fail tool-miss policy, evidence rewrite from file content,
context-budget skip, the bounded worker pool in `Engine.Run` with stats
aggregation, and session-JSONL logging with an injectable directory. All
behaviors proven against fake conversations/clients.

**Done when:** the six listed ids are each discharged by a tagged test in
`internal/engine`, and the suite is green.
