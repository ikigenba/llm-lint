# llm-lint — Plan Status

One line per pending phase in build order; this is the only home of phase
status markers. The build loop finds its next work with
`grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1` and reads only
that phase's body file, deleting the line and the file on completion.

Next phase: 28

- Phase 25 ⬜ realizes R-4JXW-H40T — live efficacy smoke for boolean-state-machine
- Phase 26 ⬜ realizes — — widen the self-lint gate config to the whole catalog (depends on Phase 25)
- Phase 27 ⬜ realizes R-H1KL-BCFY, R-H2SH-P46N — swap the verbose line to id-first ordering
