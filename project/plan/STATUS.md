# llm-lint — Plan Status

One line per pending phase in build order; this is the only home of phase
status markers. The build loop finds its next work with
`grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1` and reads only
that phase's body file, deleting the line and the file on completion.

Next phase: 11

- Phase 07 ⬜ realizes R-H8SC-8KMV R-HA08-MCDK R-HB85-0449 R-HCG1-DVUY — inline suppression filter
- Phase 08 ⬜ realizes R-HDNX-RNLN R-HEVU-5FCC R-HG3Q-J731 R-HHBM-WYTQ — reporting and full pipeline wiring
- Phase 09 ⬜ realizes R-GZ15-6EPB R-H091-K6G0 — live provider smoke tests
- Phase 10 ⬜ realizes R-HJRF-OIB4 R-HKZC-2A1T — release packaging
