# Phase 1 tests couple to the stub default-enabled `clarity` catalog

Observed while briefing Phase 03: completed Phase 1 `run` tests assume the
placeholder `internal/rules` catalog (`builtin/clarity.md`, `Select` treating
an empty enable list as "all rules") rather than D03/D09.

- `cmd/llm-lint/list_rules_test.go` (R-FX89-BIPS) expects a single
  default-enabled line containing `clarity`, `error`, `enabled`, and
  `flags unclear language`.
- `cmd/llm-lint/run_clean_test.go` (R-FR4R-EO0B) requires a positive
  inference-client call count with no config and no `--rules`.
- `cmd/llm-lint/run_error_finding_test.go` (R-FSCN-SFR0) and
  `cmd/llm-lint/run_warning_finding_test.go` (R-FTKK-67HP) inject findings
  under rule id `clarity` and likewise invoke `run` with an empty enable
  list, so they only fire if a rule is selected by default.

Out of scope for Phase 03 gather: those ids belong to D01, not this phase.
Implementing D03 (`Select` returns exactly the `enable` list; empty is not
an error) and D09 (catalog is exactly `no-sleep-in-tests`, not shipped
enabled) will go red on those tests unless a later committing turn updates
them to opt a rule in via config.
