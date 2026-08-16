# llm-lint — Plan Status

One line per pending phase in build order; this is the only home of phase
status markers. The build loop finds its next work with
`grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1` and reads only
that phase's body file, deleting the line and the file on completion.

Next phase: 18

- Phase 16 ⬜ realizes R-IVUP-YH7P, R-IX2M-C8YE, R-IYAI-Q0P3, R-IZIF-3SFS, R-J0QB-HK6H, R-J364-93NV, R-J4E0-MVEK, R-J5LX-0N59 — agent-repl model/provider/auth selection in `internal/config` (provider table, `auth`/`auth_file`, subscription credentials, `GEMINI_API_KEY`)
- Phase 17 ⬜ realizes R-J6TT-EEVY, R-J81P-S6MN — `-c`-only CLI with the `--help` model catalog and the live subscription smoke in `cmd/llm-lint`
