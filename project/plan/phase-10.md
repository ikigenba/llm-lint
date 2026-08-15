# Phase 10 — Release packaging

*Realizes design Decision 10 (Release and installation). Depends on Phase 1.*

End state: `.goreleaser.yaml`, `.github/workflows/release.yml`, `install.sh`,
and the Makefile per D10, mutually consistent on names, paths, platforms, and
triggers, with a release-consistency test in the ralph `release_test.go`
pattern; a README covering install and usage.

**Done when:** R-HJRF-OIB4 and R-HKZC-2A1T are each discharged by a tagged
test, and the suite is green.
