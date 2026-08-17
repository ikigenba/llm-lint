# llm-lint — Plan Status

One line per pending phase in build order; this is the only home of phase
status markers. The build loop finds its next work with
`grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1` and reads only
that phase's body file, deleting the line and the file on completion.

Next phase: 30

- Phase 29 ⬜ realizes R-4A0P-7AVY, R-4B8L-L2MN, R-4CGH-YUDC — render --help and bad-flag usage in agent-repl style
