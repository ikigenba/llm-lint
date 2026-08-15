# Phase 2 — Config resolution package

*Realizes design Decision 2 (Configuration resolution). Depends on Phase 1.*

End state: `internal/config` per D02 — `.llm-lint.json` walk-up discovery,
strict-JSON parsing with the flat agent-repl model keys, `-c` overrides with
`default` reset, the `gemini-3.7-flash` default, `ErrConfig`/`ErrAuth`
surfaces, and `NewConversation` building a fresh catalog-resolved agentkit
Conversation with env-key auth and attached log. `run` switches from its
Phase-1 stub to this package.

**Done when:** R-FZO2-3276, R-G0VY-GTXV, R-G23U-ULOK, R-G3BR-8DF9,
R-G5RJ-ZWWN, R-G6ZG-DONC, R-G87C-RGE1, R-G9F9-584Q are each discharged by a
tagged test in `internal/config`, and the suite is green.
