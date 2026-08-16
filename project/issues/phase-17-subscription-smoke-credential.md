# Phase 17 subscription smoke requires an unavailable credential

`cmd/llm-lint/live_subscription_auth_test.go` follows the phase brief by
skipping unless the default OpenAI subscription token at
`~/.llm-lint/openai-auth.json` exists and loads. That token is absent in the
loop environment, so the test cannot exercise the real subscription transport.

The phase done bar permits credential-gated skips only for the Google and
OpenAI API-key smokes, which makes the subscription requirement unsatisfiable
without either provisioning a real subscription token to the loop environment
or reshaping the requirement into a runnable mechanism check. Weakening the
test to pass without making a real provider call would not prove the required
end-to-end behavior.
