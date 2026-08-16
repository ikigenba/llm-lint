# Phase 26 — Widen the self-lint gate to the whole catalog

*Realizes design Decision 11 (self-lint gate) — structural. Depends on
Phase 25.*

The repo's root `.llm-lint.json` drops its `enable` allowlist so the
self-lint gate runs every catalog rule (D03's default-on behavior), keeping
only the `testdata/**` exclude. End state: the file is exactly D11's declared
content. If the keyed verify-turn lint then flags real code in this tree,
each finding is either fixed or judged a false positive and silenced with an
inline `llm-lint:ignore` comment (D07) — the gate still requires zero
findings.

**Done when:**

- `jq -e '(has("enable") | not) and .exclude == ["testdata/**"]' .llm-lint.json`
  exits 0.
- The suite is green (per design CONVENTIONS.md).
