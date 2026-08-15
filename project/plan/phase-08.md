# Phase 8 — Reporting and full pipeline wiring

*Realizes design Decision 8 (Reporting, formats, and exit computation).
Depends on Phases 2, 3, 4, 5, 6, 7.*

End state: `internal/report` per D08 — deterministic sort, text and JSON
renderers, the stats line — and `run` composing the real pipeline end to end:
config → rules → walk → cache-wrapped engine → suppression → report, with
Phase-1 stubs gone. Integration tests drive `run` with a fake engine client
through real config/rules/walk/cache/suppress on tree fixtures.

**Done when:** R-HDNX-RNLN, R-HEVU-5FCC, R-HG3Q-J731, R-HHBM-WYTQ are each
discharged by a tagged test, and the suite is green.
