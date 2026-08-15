# llm-lint

A standalone Go CLI linter whose rules are prompts and whose verdicts come
from inference: it runs a catalog of prompt-defined anti-pattern rules against
a code tree using a cheap model via
[agentkit](https://github.com/ikigenba/agentkit) and reports violations
exactly like a conventional linter — vet-style output, opt-in rules, inline
suppression, content-hash caching. Module `github.com/ikigenba/llm-lint`,
binary `llm-lint` from `cmd/llm-lint`.

## How changes are made

Changes go through the spec under `project/`, not direct edits — settle the
spec, then let the build loop realize it. The spec itself is direction-gated:
`project/**` is written only inside an operator-invoked move (the `$open-spec`
→ `$grill-me` → `$seal-spec` arc, or the build loop's completion mutations).
In any other session `project/` is read-only reference — a stale or wrong spec
is a finding to report, not a license to edit, and a settled discussion is not
direction: say what should change and wait. Edit code directly only on
explicit operator instruction. See the `$ikispec` skill for the `project/`
spec contracts and `$ralph` for the unattended build workflow.

## Layout

- `cmd/llm-lint/` — shallow composition root: a one-line `main`, the
  `version` linker seam, and the testable `run(args, in, out, errOut, getenv,
  cwd) int` every exit path flows through.
- `internal/` — all logic, run against injected deps: `config` (`.llm-lint.json`
  discovery, flat agent-repl model keys, `-c` overrides), `rules` (frontmatter
  parsing, the embedded built-in catalog under `rules/builtin/`), `walk`
  (git-based candidate enumeration), `engine` (one inference call per
  rule×file pair through the `Client` seam; tests use fakes), `cache` (verdict
  cache decorator), `suppress` (`llm-lint:ignore` filtering), `report` (text/
  JSON rendering, exit computation).
- `project/` — the spec (product/design/plan) the build loop works from.

## Tests

- Unit: `go test ./...` (or `make test`).
- Green bar (design's Conventions): `go build ./...`, `go vet ./...`, and
  `go test ./...` exit 0, and `gofmt -l .` prints nothing (excluding
  `project/`).
- Live provider smokes `t.Skip` when `GOOGLE_API_KEY` / `OPENAI_API_KEY` is
  absent; the suite must be green with no network and no credentials.

## Versioning

Versions are annotated git tags only, `vMAJOR.MINOR.PATCH` — no `VERSION`
file; `main.version` defaults to `dev` and is stamped at link time. Add the
`CHANGELOG.md` entry in the release commit itself, `git tag -a vX.Y.Z -m
"vX.Y.Z"` on `main`, then push branch and tag together with
`git push --follow-tags` — the `v*` tag push runs goreleaser via
`.github/workflows/release.yml`, publishing the archives `install.sh`
fetches. Latest is `git tag --sort=-v:refname | head -1`. The changelog
follows the agentkit convention: one `## vX.Y.Z` section per release, newest
first, flat past-tense full-sentence bullets describing what changed for a
user of the binary.
