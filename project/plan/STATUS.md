# llm-lint — Plan Status

One line per pending phase in build order; this is the only home of phase
status markers. The build loop finds its next work with
`grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1` and reads only
that phase's body file, deleting the line and the file on completion.

Next phase: 28

- Phase 27 ⬜ realizes R-H1KL-BCFY, R-H2SH-P46N — swap the verbose line to id-first ordering
