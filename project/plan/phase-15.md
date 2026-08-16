# Phase 15 — drop the model from the cache key

*Realizes design Decision 6 (Verdict cache). Structural.*

`cache.Key` loses its `model` parameter and hashes only the rule content and the
file content, matching D06's revised signature. The `CachingClient.Judge` call
site drops the empty-string argument it currently passes. This is a pure seam
cleanup: the model was already hashed as an empty string, so cache behavior is
unchanged — the code and the design simply stop carrying a parameter that must
have no effect.

**Done when:**

- `cache.Key` in `internal/cache/cache.go` is `func Key(rulePromptAndMeta,
  fileContent []byte) string` — no `model` parameter — and the string `model`
  no longer appears in `internal/cache/cache.go` (`grep -n 'model'
  internal/cache/cache.go` prints nothing).
- R-H1GX-XY6P stays green: an immediate identical re-run (same rules and files)
  makes zero calls on the fake Client and reports identical findings.
- The suite is green (`gofmt -l .` excluding `project/` empty, `go vet ./...`
  clean, `go test ./...` exits 0).
