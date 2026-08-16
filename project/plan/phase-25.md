# Phase 25 — Live efficacy smoke for boolean-state-machine

*Realizes design Decision 12 (built-in rule: boolean-state-machine — the live
smoke slice).*

A fixture tree under `testdata/` containing a blatant boolean-state-machine
violation — a Go struct with `Held` and `Sold` boolean fields and a
state-change function that writes both flags together — plus a fixture
`.llm-lint.json` enabling `boolean-state-machine`, and a live test mirroring
the shape of the existing R-GZ15-6EPB smoke: the assembled binary path (`run`
through the real composition root), default model, `t.Skip` without
`GEMINI_API_KEY`.

**Done when:**

- R-4JXW-H40T — live (Google, skips without `GEMINI_API_KEY`): the run over
  the fixture reports a finding in the fixture file and returns exit 1 —
  tagged by a genuine test.
- The suite is green (per design CONVENTIONS.md), including with no
  credentials present (the smoke skips).
