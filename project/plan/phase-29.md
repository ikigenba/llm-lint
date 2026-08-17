# Phase 29 — render --help and bad-flag usage in agent-repl style

*Realizes design Decision 1 (CLI surface and shallow composition root).*

Rework the help and usage rendering in `cmd/llm-lint` to follow agent-repl's
layout (see D1's `--help` section). `--help` prints, in fixed order: the
`usage: llm-lint [flags] [path ...]` line; a `flags:` block with one aligned
row per flag and a short description; a static `defaults:` line
(`provider=google   model=gemini-3.7-flash   auth=key`, compiled-in, never the
resolved config); an aligned `providers:` block (`<provider>  auth=key
(ENVVAR)` with `auth=sub` continuation rows); per-provider model sections
rendering each reasoning clause as `<term>={<choices>}` with the default choice
`*`-prefixed, numeric ranges as `<min>–<max>`, provider-default budgets as
`dynamic`, and no-reasoning models as `—`; and the existing `bare model=` /
`explicit provider=` footer. The bad-flag path prints the stdlib flag-error
line followed by that same single `usage:` line (no `flags:` block, no
catalog), exit 2.

This replaces the current two-line `Usage:`/`Options:` rendering, so the
already-realized bad-flag and help tests must be updated to the new output:
`R-FUSG-JZ8E` and `R-WBQM-PELO` (`unknown_flag_test.go`) now assert the
`usage: llm-lint [flags] [path ...]` line, and `R-J6TT-EEVY`
(`help_providers_catalog_test.go`) now asserts the aligned providers rows and
per-provider model headers. `--help` still exits 0, writes to stdout, needs no
network, and records zero inference calls.

**Done when:** the following ids are each covered by a genuine test in
`cmd/llm-lint/*_test.go` and the suite is green —

- R-4A0P-7AVY — `--help` opens with the exact `usage: llm-lint [flags]
  [path ...]` line and a `flags:` block with one row per flag, exit 0.
- R-4B8L-L2MN — `--help`'s `defaults:` line reads exactly `provider=google`,
  `model=gemini-3.7-flash`, `auth=key`, and is unchanged by a cwd
  `.llm-lint.json` setting a different model and by `-c model=...`, exit 0.
- R-4CGH-YUDC — a model's reasoning clause renders `<term>={<choices>}` with
  only the default choice `*`-prefixed (e.g.
  `thinking_level={low|*medium|high}`), numeric ranges as `<min>–<max>`,
  provider-default budgets as `dynamic`, and no-reasoning models as `—`.

and the updated `R-FUSG-JZ8E`, `R-WBQM-PELO`, and `R-J6TT-EEVY` tests still pass
against the new rendering; `go build ./...`, `go vet ./...`, `go test ./...`
exit 0 and `gofmt -l .` is empty (excluding `project/`).
