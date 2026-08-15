# CachingClient does not forward ContextWindow

Observed while briefing Phase 08 (reading `internal/cache` and `internal/engine` public surfaces for the pipeline-wiring dependency list).

`engine.Engine.Run` optionally type-asserts its `Client` to `interface{ ContextWindow() int64 }` before skipping oversized files. The live client (`engine.AgentkitClient`) implements `ContextWindow()`, but the wrapper this phase is told to compose around it — `cache.CachingClient` — does not, so a cache-wrapped engine never sees the budget.

Out of scope for Phase 08 gather: ContextWindow belongs to D05/D06, not D08. Filing only; the brief still lists `CachingClient` as the wrap `run` compiles against.
