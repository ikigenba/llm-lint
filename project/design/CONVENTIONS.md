# llm-lint — Design Conventions

- **Language**: Go 1.26. Module path `github.com/ikigenba/llm-lint`.
- **Layout**: `cmd/llm-lint/` (shallow main + `run`), `internal/<pkg>/` for
  everything else. The binary is `llm-lint`.
- **Build/typecheck command**: `go build ./...`
- **Test command**: `go test ./...`
- **The suite is green** when `gofmt -l .` prints nothing (excluding
  `project/`), `go vet ./...` is clean, and `go test ./...` exits 0. Live
  integration tests `t.Skip` when their provider env key is absent, so green
  does not require credentials.
- **Test-file glob**: `*_test.go`. Requirement-id tags appear verbatim in test
  files, in a comment on (or subtest name of) the test that discharges them.
- **Dependencies**: `github.com/ikigenba/agentkit` (client, catalog,
  providers) and `github.com/bmatcuk/doublestar/v4` (globs). Nothing else
  without a new Decision.
- **Exit-code taxonomy** (used verbatim everywhere): `0` clean (or
  warnings-only), `1` at least one error-severity finding, `2` usage/config/
  rule-file error, `3` operational failure (provider, auth, network, I/O).
- **Formatting**: `gofmt`. Errors are wrapped with `%w` and package-prefixed
  messages (`config: …`, `rules: …`).
- **Time and environment**: injected at package boundaries (`getenv
  func(string) string`, `now func() time.Time`); nothing below `cmd/` reads
  `os.Getenv`, `os.Args`, or the ambient clock directly.
