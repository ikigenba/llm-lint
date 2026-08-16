# AgentkitClient.Context is never populated by run()

Observed while briefing Phase 18 (reading `cmd/llm-lint` composition-root wiring against D05's oversize-skip path).

`engine.AgentkitClient` exposes `Context int64` and `ContextWindow() int64`, and `Engine.Run` skips a pair only when `ContextWindow()` returns a positive budget that the bytes/3 estimate exceeds. `newClient` in `cmd/llm-lint/main.go` constructs `AgentkitClient` with `NewConversation` and `Warn` only — `Context` stays 0 — so even a client that is not cache-wrapped never takes the skip path in production.

Out of scope for Phase 18 gather: R-GWLC-EV7X already landed against a fake `budgetClient`, and this phase's skip-trace id (R-H40E-2VXC) is likewise an engine-level test. Filing only; the brief still treats skip recording as an `Engine` concern gated on `ContextWindow()`.
