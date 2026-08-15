# Phase 1 — Module scaffold, shallow main, and the run composition root

*Realizes design Decision 1 (CLI surface and shallow composition root).*

End state: a buildable module `github.com/ikigenba/llm-lint` (go 1.26,
agentkit and doublestar in `go.mod`) with `cmd/llm-lint/main.go` holding the
one-line `main`, the `version` variable, and `run(args, in, out, errOut,
getenv, cwd) int` parsing every flag from D01 and returning taxonomy exit
codes. Downstream seams (config, rules, walker, engine, cache, suppression,
report) exist as the minimal interfaces/stubs `run` composes, sufficient for
fakes to drive every D01 behavior: clean exit, finding exits by severity,
usage errors, `--version`, `--list-rules` without inference, and operational
failure mapping. `.gitignore` for Go artifacts included.

**Done when:** R-FR4R-EO0B, R-FSCN-SFR0, R-FTKK-67HP, R-FUSG-JZ8E,
R-FW0C-XQZ3, R-FX89-BIPS, R-FYG5-PAGH are each discharged by a tagged test in
`cmd/llm-lint`, and the suite is green (`gofmt -l .` empty, `go vet ./...`
clean, `go test ./...` exit 0).
