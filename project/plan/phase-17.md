# Phase 17 — `-c`-only CLI with the `--help` model catalog and the live subscription smoke

*Realizes design Decision 1 (CLI surface) and 5 (Inference engine), the CLI/
help and live-test slices. Depends on Phase 16.*

Finish the user surface in `cmd/llm-lint`: make `-c` the sole channel for
agentkit settings and teach `--help` to print the model catalog, then align the
live smokes with the new provider/auth defaults and add the subscription smoke.

Observable end state:

- The `--model` flag is gone; model configuration flows only through `-c`. A
  `--model …` invocation is now an unknown flag (usage to stderr, exit 2, via
  the existing R-FUSG-JZ8E path). Any test that drove `--model` is updated to
  `-c model=…`.
- `--help` prints, after llm-lint's own flags, a `providers:` block built from
  D02's `Providers()`/`ProviderInfo` (each provider's `auth=key (ENVVAR)` line,
  plus an `auth=sub (auth_file=~/.llm-lint/<provider>-auth.json)` line for
  openai and x-ai), then each provider's catalogued models
  (`agentkit/catalog.ListCurated`) with reasoning clauses marking the default,
  then the bare-`model=` vs. explicit-`provider=` footer. It exits 0 and makes
  no inference call.
- `Config.Warnings` are printed once to stderr before linting.
- The Google live smoke gates on `GEMINI_API_KEY`; the OpenAI live smoke runs
  `-c model=gpt-5.6-luna -c auth=key`; a new subscription smoke exercises
  `-c auth=sub` end-to-end.

**Done when** the following are covered by genuine tests and the suite is green
(`go build ./...`, `go vet ./...`, `go test ./...` exit 0; `gofmt -l .`
excluding `project/` prints nothing):

- R-J6TT-EEVY — a `--help` test asserts the output contains the `providers:`
  block naming each provider with its env key (google → `GEMINI_API_KEY`), the
  `auth=sub` lines for openai and x-ai, and catalogued model ids (e.g.
  `gemini-3.7-flash`) with reasoning clauses; exit 0 and zero inference calls
  recorded.
- R-J81P-S6MN — a live subscription smoke (`cmd/llm-lint`, through the real
  composition root) that `t.Skip`s unless the sub-capable provider's token file
  `~/.llm-lint/<provider>-auth.json` exists and loads; when present, judging the
  blatant-sleep fixture with `-c auth=sub` reports the finding and returns
  exit 1, proving the real `Subscription` transport.
- R-GZ15-6EPB and R-H091-K6G0 — their live tests are updated to the new env var
  (`GEMINI_API_KEY`) and the explicit `-c auth=key` on the OpenAI path, still
  `t.Skip`ping without their credentials so credential-free green is unchanged.
