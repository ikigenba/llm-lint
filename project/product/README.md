# llm-lint — Product

**Authority: intent.** This document owns the why: the problem, who llm-lint
is for, its scope, and what we promise the user, stated in outcome terms. It
does not own mechanism, exact formats, exit codes, or test assertions — those
belong to design. Where behavior appears here it is a promise; design states
the exact, checkable proof of each promise.

## Problem

Some code defects are anti-patterns no static analyzer can catch because
recognizing them takes judgment, not parsing. The canonical example: a sleep
in a test that papers over a timing or synchronization issue. A grep for
`time.Sleep` drowns in legitimate uses; telling "this sleep waits out a race"
from "this sleep is the behavior under test" requires reading the code the way
a reviewer would. Today that judgment only happens when a human reviewer
happens to look, so these defects accumulate silently — especially in
codebases built by autonomous loops, where no human reads most diffs.

## Purpose

llm-lint is a linter whose rules are prompts and whose verdicts come from
inference. It runs a catalog of prompt-defined anti-pattern rules against a
code tree using a cheap language model and reports violations exactly the way
a conventional linter would. In every respect other than how verdicts are
reached — invocation, configuration, output, suppression, CI behavior — it
works like the linters developers already know.

## Users

- **Developers** running it locally against their working tree, and wiring it
  into CI alongside their other linters.
- **Autonomous build loops** (the ikigenba ralph/verify gates) running it as a
  mechanical quality gate and acting on its findings, including suppressing a
  finding they have judged a false positive.

## Scope

llm-lint lints files against prompt-defined rules and reports findings; it
does nothing else. It ships a built-in rule catalog (starting with a single
rule: sleeps in tests papering over timing issues) and runs project-supplied
rule files alongside it. Every rule available to a run is active by default; a
project narrows the run to a chosen subset or silences individual rules by name
in its config. It never modifies the code it lints, never manages adjudication or override
state beyond honoring inline suppression comments, and never runs as a daemon
or editor service — it is a batch CLI. Judgment quality is bounded by the
configured model; the tool's job is to make running that judgment cheap,
repeatable, and conventional, not to guarantee a perfect verdict.

## Contractual constants

- The binary and tool name is `llm-lint`.
- The project configuration file is named `.llm-lint.json`.
- The inline suppression marker is `llm-lint:ignore`.
- The default model is `gemini-3.7-flash`.
- The module and release home is `github.com/ikigenba/llm-lint`.

## What we promise (user-facing behavior)

- Running `llm-lint` in a configured project examines the tree and prints one
  line per violation, naming the file, line, rule, and what the model found —
  and prints nothing on a clean run:

      $ llm-lint ./...
      internal/poller/poller_test.go:41: fixed 100ms sleep waits for the goroutine instead of synchronizing (no-sleep-in-tests)

- CI can gate on it exactly as it gates on any linter: success and failure are
  distinguishable, and a broken tool (no credentials, provider outage) is
  distinguishable from a failing lint, so an outage never reads as either a
  pass or a code defect.
- Every built-in rule runs by default; a project can restrict a run to a
  chosen subset of rules or silence individual rules by name in
  `.llm-lint.json`. Projects add their own rules as committed prompt files that
  behave identically to built-in ones and run alongside them, and can list
  every rule available to a run and see which are active.
- A finding judged acceptable is silenced by an ordinary inline comment
  (`llm-lint:ignore`), committed next to the code, exactly like other linters'
  suppression comments — this is also how build loops override a false
  positive.
- Re-running on an unchanged tree is fast and spends nothing: verdicts are
  remembered per file and rule content, and only what changed is re-judged.
- With a provider API key in the environment it works out of the box on the
  default model. Any of its supported models can be selected by name — using
  that model's default provider, or a provider named explicitly — and its
  generation settings are configurable per project and per invocation. All of
  this model configuration is set the same way, through repeatable config
  overrides, with no dedicated per-setting flags.
- Where a provider offers a subscription, a run can authenticate against that
  subscription instead of a metered API key — the cheaper path — by pointing
  it at a subscription token file, chosen entirely in configuration.
- Asking for help lists every model available to a run, grouped by provider,
  with each provider's authentication options, so a user can see what they can
  select without leaving the terminal.
- Machine-readable output is available for tooling, carrying the exact
  offending source line so consumers can build stable identities for findings.
- Every run leaves a complete forensic log of its model exchanges for
  after-the-fact inspection, and can report what a run cost on request.
- It installs the way its siblings do: a released binary fetched by a shell
  installer, runnable on Linux and macOS.

## Success criteria (outcomes)

- A project with an obvious sleep-to-dodge-a-race in a test and a valid API
  key gets that sleep flagged with the right file and line by a real model run.
- With no configuration every built-in rule is active; a project can restrict a
  run to a chosen subset or disable a rule by name, and see the result in the
  rule listing.
- A clean project produces no output and CI treats the run as a success;
  introducing a violation flips the same CI job to failure; removing the API
  key produces a tool-failure outcome CI can tell apart from both.
- A project-authored rule file, once configured, produces findings with
  the same output, suppression, and caching behavior as the built-in rule.
- Adding an inline `llm-lint:ignore` comment to a flagged line makes the
  finding disappear from the next run without changing anything else.
- Re-running immediately after a completed run on an unchanged tree completes
  without any new model calls and reports the same findings.
- A user can list every available rule and see which are enabled.
- The released artifact installs via the installer script onto a clean Linux
  or macOS machine and reports its release version.
- Running `llm-lint --help` lists every selectable model grouped by provider,
  with each provider's authentication options.
- A model runs against its default provider when only its name is given and
  against a chosen provider when one is set, authenticated by an env key or,
  where the provider offers it, a cheaper subscription token — all configured
  without code changes.
