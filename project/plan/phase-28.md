# Phase 28 — print the bad-flag usage block exactly once

*Realizes design Decision 1 (CLI surface and shallow composition root).*

On a flag-parse failure `run` currently emits the usage block twice: the
stdlib `flag.ContinueOnError` handler already prints its error line and renders
the flag set's `Usage` (one block), and `run` then renders the block a second
time on the parse-error path. Remove the redundant render in
`cmd/llm-lint/main.go` so a bad flag produces the stdlib flag-error line
followed by a single usage block on stderr, exit 2. `R-FUSG-JZ8E` (usage
present, exit 2) stays green; the new behavior is that the block appears exactly
once.

**Done when:** a test tagged `R-WBQM-PELO` drives `run` with an unknown flag
and asserts the usage block appears exactly once on stderr (the substring
`Usage: llm-lint` occurs a single time), preceded by the stdlib flag-error
line, with return code 2; the existing `R-FUSG-JZ8E` assertion still passes; and
the suite is green (`go build ./...`, `go vet ./...`, `go test ./...` exit 0,
`gofmt -l .` empty excluding `project/`).
