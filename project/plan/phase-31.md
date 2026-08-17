# Phase 31 — bump the agentkit dependency to v0.22.0

*Realizes design Decision — (structural: no new Verification ids).*

Update the module's dependency on `github.com/ikigenba/agentkit` from v0.20.0
to v0.22.0 in `go.mod` and `go.sum`. v0.22.0 is a minor, additive release: its
only exported-API change is the two new `catalog.VendorID` constants
(`VendorNVIDIA`, `VendorQwen`), and it adds three OpenRouter-served catalog
models (`nemotron-3.5-lightning`, `qwen3.8-max`, `qwen3.8-27b`) plus
`claude-opus-5` from v0.21.0. llm-lint needs no source change: these are pure
catalog data that its existing catalog-driven config resolution and `--help`
enumeration already consume, so after the bump a bare
`-c model=qwen3.8-max` (or any new model) resolves to its OpenRouter provider
through the same code path every curated model uses, and the new models appear
in `--help` automatically.

End state: `go.mod` requires `github.com/ikigenba/agentkit v0.22.0`, `go.sum`
carries its checksums, and the suite is green with no source edits (the
existing help tests use substring assertions, not exact model-set matches, so
the added models do not break them).

**Done when:**

- `go.mod` shows `github.com/ikigenba/agentkit v0.22.0` (e.g.
  `grep -qE '^\s*github\.com/ikigenba/agentkit v0\.22\.0\b' go.mod`).
- The green bar holds: `go build ./...`, `go vet ./...`, and `go test ./...`
  exit 0, and `gofmt -l .` prints nothing (excluding `project/`).
