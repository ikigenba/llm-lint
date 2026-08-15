# llm-lint — Plan Status

One line per pending phase in build order; this is the only home of phase
status markers. The build loop finds its next work with
`grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1` and reads only
that phase's body file, deleting the line and the file on completion.

Next phase: 11

- Phase 01 ⬜ realizes R-FR4R-EO0B R-FSCN-SFR0 R-FTKK-67HP R-FUSG-JZ8E R-FW0C-XQZ3 R-FX89-BIPS R-FYG5-PAGH — module scaffold, shallow main, and the run composition root
- Phase 02 ⬜ realizes R-FZO2-3276 R-G0VY-GTXV R-G23U-ULOK R-G3BR-8DF9 R-G5RJ-ZWWN R-G6ZG-DONC R-G87C-RGE1 R-G9F9-584Q — config resolution package
- Phase 03 ⬜ realizes R-GAN5-IZVF R-GBV1-WRM4 R-GD2Y-AJCT R-GEAU-OB3I R-GFIR-22U7 R-GGQN-FUKW R-GHYJ-TMBL R-HIJJ-AQKF — rules package with embedded built-in catalog
- Phase 04 ⬜ realizes R-GJ6G-7E2A R-GKEC-L5SZ R-GLM8-YXJO R-GO21-QH12 R-GP9Y-48RR R-GQHU-I0IG — candidate file walker
- Phase 05 ⬜ realizes R-GRPQ-VS95 R-GSXN-9JZU R-GU5J-NBQJ R-GVDG-13H8 R-GWLC-EV7X R-GXT8-SMYM — inference engine against fakes
- Phase 06 ⬜ realizes R-H1GX-XY6P R-H2OU-BPXE R-H3WQ-PHO3 R-H54N-39ES R-H6CJ-H15H — verdict cache decorator
- Phase 07 ⬜ realizes R-H8SC-8KMV R-HA08-MCDK R-HB85-0449 R-HCG1-DVUY — inline suppression filter
- Phase 08 ⬜ realizes R-HDNX-RNLN R-HEVU-5FCC R-HG3Q-J731 R-HHBM-WYTQ — reporting and full pipeline wiring
- Phase 09 ⬜ realizes R-GZ15-6EPB R-H091-K6G0 — live provider smoke tests
- Phase 10 ⬜ realizes R-HJRF-OIB4 R-HKZC-2A1T — release packaging
