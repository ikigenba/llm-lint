# llm-lint — Plan Status

One line per pending phase in build order; this is the only home of phase
status markers. The build loop finds its next work with
`grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1` and reads only
that phase's body file, deleting the line and the file on completion.

Next phase: 15

- Phase 13 ⬜ realizes R-2J16-KM75, R-2P4O-HGWM — cmd: wire default-on, `disable`, and the `-V` alias
- Phase 14 ⬜ realizes — ship the self-lint config (`.llm-lint.json`)
