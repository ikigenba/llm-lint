# llm-lint — Plan Status

One line per pending phase in build order; this is the only home of phase
status markers. The build loop finds its next work with
`grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1` and reads only
that phase's body file, deleting the line and the file on completion.

Next phase: 25

- Phase 24 ⬜ realizes R-4HI3-PKJF, R-GBV1-WRM4, R-4IQ0-3CA4, R-HIJJ-AQKF — default code globs and the boolean-state-machine built-in rule
