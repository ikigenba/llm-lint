# llm-lint — Design Index

Each Decision maps to its `DNN.md`; every Verification id maps to its
Decision and file; id lookup is a grep against this index. Regenerated
whenever a Decision is added or its Verification ids change.

## Decisions

- D1 → `D01.md` — CLI surface and shallow composition root — R-FR4R-EO0B,
  R-FSCN-SFR0, R-FTKK-67HP, R-FUSG-JZ8E, R-FW0C-XQZ3, R-FX89-BIPS,
  R-FYG5-PAGH
- D2 → `D02.md` — Configuration resolution — R-FZO2-3276, R-G0VY-GTXV,
  R-G23U-ULOK, R-G3BR-8DF9, R-G5RJ-ZWWN, R-G6ZG-DONC, R-G87C-RGE1,
  R-G9F9-584Q
- D3 → `D03.md` — Rule model and catalog — R-GAN5-IZVF, R-GBV1-WRM4,
  R-GD2Y-AJCT, R-GEAU-OB3I, R-GFIR-22U7, R-GGQN-FUKW, R-GHYJ-TMBL
- D4 → `D04.md` — Candidate file walking — R-GJ6G-7E2A, R-GKEC-L5SZ,
  R-GLM8-YXJO, R-GO21-QH12, R-GP9Y-48RR, R-GQHU-I0IG
- D5 → `D05.md` — Inference engine — R-GRPQ-VS95, R-GSXN-9JZU, R-GU5J-NBQJ,
  R-GVDG-13H8, R-GWLC-EV7X, R-GXT8-SMYM, R-GZ15-6EPB, R-H091-K6G0
- D6 → `D06.md` — Verdict cache — R-H1GX-XY6P, R-H2OU-BPXE, R-H3WQ-PHO3,
  R-H54N-39ES, R-H6CJ-H15H
- D7 → `D07.md` — Inline suppression — R-H8SC-8KMV, R-HA08-MCDK,
  R-HB85-0449, R-HCG1-DVUY
- D8 → `D08.md` — Reporting, formats, and exit computation — R-HDNX-RNLN,
  R-HEVU-5FCC, R-HG3Q-J731, R-HHBM-WYTQ
- D9 → `D09.md` — Built-in rule: no-sleep-in-tests — R-HIJJ-AQKF
- D10 → `D10.md` — Release and installation — R-HJRF-OIB4, R-HKZC-2A1T

## Verification ids → Decision

- R-FR4R-EO0B — D1 (`D01.md`)
- R-FSCN-SFR0 — D1 (`D01.md`)
- R-FTKK-67HP — D1 (`D01.md`)
- R-FUSG-JZ8E — D1 (`D01.md`)
- R-FW0C-XQZ3 — D1 (`D01.md`)
- R-FX89-BIPS — D1 (`D01.md`)
- R-FYG5-PAGH — D1 (`D01.md`)
- R-FZO2-3276 — D2 (`D02.md`)
- R-G0VY-GTXV — D2 (`D02.md`)
- R-G23U-ULOK — D2 (`D02.md`)
- R-G3BR-8DF9 — D2 (`D02.md`)
- R-G5RJ-ZWWN — D2 (`D02.md`)
- R-G6ZG-DONC — D2 (`D02.md`)
- R-G87C-RGE1 — D2 (`D02.md`)
- R-G9F9-584Q — D2 (`D02.md`)
- R-GAN5-IZVF — D3 (`D03.md`)
- R-GBV1-WRM4 — D3 (`D03.md`)
- R-GD2Y-AJCT — D3 (`D03.md`)
- R-GEAU-OB3I — D3 (`D03.md`)
- R-GFIR-22U7 — D3 (`D03.md`)
- R-GGQN-FUKW — D3 (`D03.md`)
- R-GHYJ-TMBL — D3 (`D03.md`)
- R-GJ6G-7E2A — D4 (`D04.md`)
- R-GKEC-L5SZ — D4 (`D04.md`)
- R-GLM8-YXJO — D4 (`D04.md`)
- R-GO21-QH12 — D4 (`D04.md`)
- R-GP9Y-48RR — D4 (`D04.md`)
- R-GQHU-I0IG — D4 (`D04.md`)
- R-GRPQ-VS95 — D5 (`D05.md`)
- R-GSXN-9JZU — D5 (`D05.md`)
- R-GU5J-NBQJ — D5 (`D05.md`)
- R-GVDG-13H8 — D5 (`D05.md`)
- R-GWLC-EV7X — D5 (`D05.md`)
- R-GXT8-SMYM — D5 (`D05.md`)
- R-GZ15-6EPB — D5 (`D05.md`)
- R-H091-K6G0 — D5 (`D05.md`)
- R-H1GX-XY6P — D6 (`D06.md`)
- R-H2OU-BPXE — D6 (`D06.md`)
- R-H3WQ-PHO3 — D6 (`D06.md`)
- R-H54N-39ES — D6 (`D06.md`)
- R-H6CJ-H15H — D6 (`D06.md`)
- R-H8SC-8KMV — D7 (`D07.md`)
- R-HA08-MCDK — D7 (`D07.md`)
- R-HB85-0449 — D7 (`D07.md`)
- R-HCG1-DVUY — D7 (`D07.md`)
- R-HDNX-RNLN — D8 (`D08.md`)
- R-HEVU-5FCC — D8 (`D08.md`)
- R-HG3Q-J731 — D8 (`D08.md`)
- R-HHBM-WYTQ — D8 (`D08.md`)
- R-HIJJ-AQKF — D9 (`D09.md`)
- R-HJRF-OIB4 — D10 (`D10.md`)
- R-HKZC-2A1T — D10 (`D10.md`)

## Success criteria → ids

- "A project with an obvious sleep-to-dodge-a-race … gets that sleep flagged
  … by a real model run" → R-GZ15-6EPB, R-H091-K6G0
- "A clean project produces no output and CI treats the run as a success;
  introducing a violation flips the same CI job to failure; removing the API
  key produces a tool-failure outcome CI can tell apart from both" →
  R-FR4R-EO0B, R-FSCN-SFR0, R-G87C-RGE1, R-FYG5-PAGH
- "A project-authored rule file, opted in via config, produces findings with
  the same output, suppression, and caching behavior as the built-in rule" →
  R-GEAU-OB3I, R-GD2Y-AJCT
- "Adding an inline `llm-lint:ignore` comment to a flagged line makes the
  finding disappear from the next run" → R-H8SC-8KMV, R-HB85-0449
- "Re-running immediately after a completed run on an unchanged tree
  completes without any new model calls and reports the same findings" →
  R-H1GX-XY6P
- "A user can list every available rule and see which are enabled" →
  R-FX89-BIPS
- "The released artifact installs via the installer script onto a clean
  Linux or macOS machine and reports its release version" → R-HJRF-OIB4,
  R-HKZC-2A1T, R-FW0C-XQZ3
