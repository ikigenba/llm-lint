# Phase 14 — ship the self-lint config

*Realizes design Decision 11 (self-lint gate). Structural.*

Create the root `.llm-lint.json` that turns this repo into a dogfood target for
its own linter: it enables the sole built-in rule and excludes the deliberately
anti-pattern fixtures under `testdata/` from every run. No Go code changes; this
is the one config file D11 ships. The keyed lint gate itself lives in the
regenerated verify prompt (a separate, operator-invoked generator step) and is
not built here.

Content:

```json
{
  "enable": ["no-sleep-in-tests"],
  "exclude": ["testdata/**"]
}
```

**Done when:**

- `.llm-lint.json` exists at the repo root and is valid JSON whose `enable`
  equals `["no-sleep-in-tests"]` and whose `exclude` equals `["testdata/**"]`.
  Deterministic check:
  `jq -e '.enable == ["no-sleep-in-tests"] and .exclude == ["testdata/**"]' .llm-lint.json`
  exits 0.
- The suite is green (`gofmt -l .` excluding `project/` empty, `go vet ./...`
  clean, `go test ./...` exits 0) — the new file changes no Go behavior.
