# Engine.Stats token and cost fields are never filled

Observed while briefing Phase 08 (reading `internal/engine` public signatures for D08's StatsLine).

`engine.Stats` exports `InputTokens`, `OutputTokens`, and `CostUSD`, and Phase 08's R-HG3Q-J731 requires the `--stats` line to include token and cost figures from the run. `Engine.Run` only populates `Rules`, `Files`, `Pairs`, and `Calls`; the three usage fields stay zero. `AgentkitClient.Judge` also does not return usage.

Out of scope for Phase 08 gather: populating those fields is engine work (D05), not reporting. The brief still has StatsLine consume the exported Stats fields as D08 specifies.
