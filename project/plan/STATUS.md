# llm-lint — Plan Status

One line per pending phase in build order; this is the only home of phase
status markers. The build loop finds its next work with
`grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1` and reads only
that phase's body file, deleting the line and the file on completion.

Next phase: 11

- Phase 05 ⬜ realizes R-GRPQ-VS95 R-GSXN-9JZU R-GU5J-NBQJ R-GVDG-13H8 R-GWLC-EV7X R-GXT8-SMYM — inference engine against fakes
- Phase 06 ⬜ realizes R-H1GX-XY6P R-H2OU-BPXE R-H3WQ-PHO3 R-H54N-39ES R-H6CJ-H15H — verdict cache decorator
- Phase 07 ⬜ realizes R-H8SC-8KMV R-HA08-MCDK R-HB85-0449 R-HCG1-DVUY — inline suppression filter
- Phase 08 ⬜ realizes R-HDNX-RNLN R-HEVU-5FCC R-HG3Q-J731 R-HHBM-WYTQ — reporting and full pipeline wiring
- Phase 09 ⬜ realizes R-GZ15-6EPB R-H091-K6G0 — live provider smoke tests
- Phase 10 ⬜ realizes R-HJRF-OIB4 R-HKZC-2A1T — release packaging
