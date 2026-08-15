# Phase 11 — config: parse the `disable` opt-out list

*Realizes design Decision 2 (Configuration resolution).*

`internal/config` gains a `Disable []string` field on `Config`, populated from
a new top-level `disable` JSON string array in `.llm-lint.json` (file order).
Absent `disable` leaves it empty; a `disable` value that is not an array of
strings is a config error naming the key, the same way other malformed keys
are. The existing absent-config default is unchanged: `Enable` and `Disable`
are both empty when no file is found.

**Done when:**

- R-2K92-YDXU — a `disable` JSON string array populates `Config.Disable` in
  file order; a non-array or non-string-element `disable` value returns
  `ErrConfig` naming the key. Covered by a genuine `internal/config` test.
- R-G0VY-GTXV still green: absent config yields empty `Enable` and `Disable`
  and the default model.
- The suite is green (`gofmt -l .` excluding `project/` empty, `go vet ./...`
  clean, `go test ./...` exits 0).
