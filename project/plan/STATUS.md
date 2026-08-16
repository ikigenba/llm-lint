# llm-lint — Plan Status

One line per pending phase in build order; this is the only home of phase
status markers. The build loop finds its next work with
`grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1` and reads only
that phase's body file, deleting the line and the file on completion.

Next phase: 20

- Phase 19 ⬜ realizes R-ETNZ-VY1Q, R-HG3Q-J731, R-H7O3-875F, R-GZ15-6EPB — reshape the Client seam for real token/cost usage; retire the oversize-skip path and the OpFailure typed error
