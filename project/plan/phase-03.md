# Phase 3 — Rules package with embedded built-in catalog

*Realizes design Decisions 3 (Rule model and catalog) and 9 (Built-in rule:
no-sleep-in-tests). Depends on Phase 1.*

End state: `internal/rules` per D03 — frontmatter `Parse`, `go:embed`
built-ins, local `Load`, `Select`, and the `ErrRule` surface — plus the
authored `internal/rules/builtin/no-sleep-in-tests.md` per D09 (frontmatter
as specified; prompt body written and readable, with one flagged and one
spared example). `run`'s rule loading and `--list-rules` use this package.

**Done when:** R-GAN5-IZVF, R-GBV1-WRM4, R-GD2Y-AJCT, R-GEAU-OB3I,
R-GFIR-22U7, R-GGQN-FUKW, R-GHYJ-TMBL, R-HIJJ-AQKF are each discharged by a
tagged test, and the suite is green.
