# llm-lint — Plan Status

One line per pending phase in build order; this is the only home of phase
status markers. The build loop finds its next work with
`grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1` and reads only
that phase's body file, deleting the line and the file on completion.

Next phase: 18

- Phase 17 ⬜ realizes R-J6TT-EEVY, R-J81P-S6MN — `-c`-only CLI with the `--help` model catalog and the live subscription smoke in `cmd/llm-lint`
