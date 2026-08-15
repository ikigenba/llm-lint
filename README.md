# llm-lint

`llm-lint` checks files with configurable LLM-backed lint rules.

## Install

Install the latest release into `$HOME/.local/bin`:

```sh
curl -fsSL https://raw.githubusercontent.com/ikigenba/llm-lint/main/install.sh | sh
```

Set `LLM_LINT_VERSION` to install a particular release tag, `BINDIR` to choose
the destination directly, or `PREFIX` to use `$PREFIX/bin`:

```sh
curl -fsSL https://raw.githubusercontent.com/ikigenba/llm-lint/main/install.sh | LLM_LINT_VERSION=v1.0.0 BINDIR=/usr/local/bin sh
```

To build and install from a checkout instead:

```sh
make install PREFIX="$HOME/.local"
```

## Usage

Run the linter against one or more paths:

```sh
llm-lint [options] [path...]
```

Use `llm-lint --help` to see all options, `llm-lint --list-rules` to list
available rules, and `llm-lint --version` to print the installed version.
