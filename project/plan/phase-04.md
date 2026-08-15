# Phase 4 — Candidate file walker

*Realizes design Decision 4 (Candidate file walking). Depends on Phase 1.*

End state: `internal/walk` per D04 — git-based enumeration with injected
`RunGit`, the non-repo `WalkDir` fallback, NUL-sniff binary exclusion, global
exclude globs, explicit-file inclusion, and per-rule `Candidates` with
doublestar semantics. `run` uses it to build the walked set.

**Done when:** R-GJ6G-7E2A, R-GKEC-L5SZ, R-GLM8-YXJO, R-GO21-QH12,
R-GP9Y-48RR, R-GQHU-I0IG are each discharged by a tagged test in
`internal/walk` (git-repo fixtures created under `t.TempDir()`), and the
suite is green.
