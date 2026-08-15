# Phase 7 — Inline suppression filter

*Realizes design Decision 7 (Inline suppression). Depends on Phase 5.*

End state: `internal/suppress` per D07 — marker recognition in any line,
rule-id lists, bare-marker suppress-all, comment-only-line next-line
targeting, and the I/O error surface.

**Done when:** R-H8SC-8KMV, R-HA08-MCDK, R-HB85-0449, R-HCG1-DVUY are each
discharged by a tagged test in `internal/suppress`, and the suite is green.
