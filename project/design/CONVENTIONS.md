# llm-lint — Design Conventions

- **Language**: Go 1.26. Module path `github.com/ikigenba/llm-lint`.
- **Layout**: `cmd/llm-lint/` (shallow main + `run`), `internal/<pkg>/` for
  everything else. The binary is `llm-lint`.
- **Build/typecheck command**: `go build ./...`
- **Test command**: `go test ./...`
- **The suite is green** when `gofmt -l .` prints nothing (excluding
  `project/`), `go vet ./...` is clean, and `go test ./...` exits 0. Live
  integration tests `t.Skip` when their required credential is absent — a
  provider env key, or a subscription token file for the sub-auth smoke — so
  green does not require credentials.
- **Self-lint gate** (D11): beyond the suite, the build loop's verify turn runs
  the installed `llm-lint` binary over the tree and requires zero findings. It
  calls a model, so this gate runs **only when a provider API key is present**
  in the environment and is **skipped otherwise** — exactly like the live
  integration tests, so credential-free green is unchanged. The repo ships a
  root `.llm-lint.json` (`exclude: ["testdata/**"]`, no `enable` key, so the
  whole catalog is active) so the gate lints real project code with every
  built-in rule and never the deliberately anti-pattern fixtures under
  `testdata/`.
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
