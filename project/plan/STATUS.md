# llm-lint — Plan Status

One line per pending phase in build order; this is the only home of phase
status markers. The build loop finds its next work with
`grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1` and reads only
that phase's body file, deleting the line and the file on completion.

Next phase: 19

- Phase 18 ⬜ realizes R-H1KL-BCFY, R-H2SH-P46N, R-H40E-2VXC, R-H58A-GNO1, R-H6G6-UFEQ, R-H7O3-875F — `--verbose` per-pair audit trace across the CLI, engine, cache, and reporter
