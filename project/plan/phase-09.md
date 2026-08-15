# Phase 9 — Live provider smoke tests

*Realizes design Decision 5 (Inference engine), slice R-GZ15-6EPB and
R-H091-K6G0. Depends on Phase 8.*

End state: guarded integration tests (suite convention: `t.Skip` when
`GOOGLE_API_KEY` / `OPENAI_API_KEY` is absent) that drive `run` through the
real composition root — real provider, real model, real catalog — against a
committed fixture tree containing a blatant sleep-past-a-race Go test with
`no-sleep-in-tests` enabled: `gemini-3.7-flash` for R-GZ15-6EPB,
`-c model=gpt-5.6-luna` for R-H091-K6G0, each asserting a finding in the
fixture file and exit 1.

**Done when:** both ids are discharged by tagged tests that pass with the
key present and skip without it, and the suite is green.
