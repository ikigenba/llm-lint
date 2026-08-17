# Phase 30 — blank-line section separators and defaults-as-header in --help

*Realizes design Decision 1 (CLI surface and shallow composition root).*

Refine the `--help` rendering in `cmd/llm-lint` to match agent-repl's spacing
(see D1's `--help` section). Insert exactly one blank line between consecutive
sections — after the `usage:` line, after the `flags:` block, after the
`defaults:` block, after the `providers:` block, and between each consecutive
per-provider model section — so no section runs together with the next.
Separately, render the `defaults:` section as a header line with its value
indented on the following line (`defaults:` then `  provider=google
model=gemini-3.7-flash   auth=key`) rather than inline on one line, matching the
header-then-rows shape the other sections already use. The bad-flag path is
unchanged: it still prints only the single `usage:` line with no sections and no
blank-line separators.

This changes the already-realized `R-4B8L-L2MN` test (defaults now a header +
indented value, not one inline line), which must be updated to the new form.

**Done when:**

- R-TIL4-XZAL — a test on `--help` output asserts exactly one blank line
  follows the `usage:` line, the `flags:` block, the `defaults:` block, and the
  `providers:` block, and separates each consecutive per-provider model section
  (no run-together sections, no two-or-more-blank-line gaps).
- the updated `R-4B8L-L2MN` test asserts the `defaults:` header and its
  indented value line (`provider=google`, `model=gemini-3.7-flash`, `auth=key`),
  still invariant to `.llm-lint.json` and `-c`.

and the suite is green: `go build ./...`, `go vet ./...`, `go test ./...` exit 0
and `gofmt -l .` is empty (excluding `project/`).
