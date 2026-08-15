#!/bin/sh
set -eu

REPO=ikigenba/llm-lint
BINARY=llm-lint
LLM_LINT_VERSION=${LLM_LINT_VERSION:-latest}
BINDIR=${BINDIR:-${PREFIX:-$HOME/.local}/bin}

case "$(uname -s)" in
    Linux) os=linux ;;
    Darwin) os=darwin ;;
    *) echo "llm-lint: unsupported operating system: $(uname -s)" >&2; exit 1 ;;
esac

case "$(uname -m)" in
    x86_64 | amd64) arch=amd64 ;;
    arm64 | aarch64) arch=arm64 ;;
    *) echo "llm-lint: unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

asset="${BINARY}_${os}_${arch}.tar.gz"
if [ "$LLM_LINT_VERSION" = latest ]; then
    url="https://github.com/${REPO}/releases/latest/download/${asset}"
else
    url="https://github.com/${REPO}/releases/download/${LLM_LINT_VERSION}/${asset}"
fi

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' 0 HUP INT TERM

curl -fsSL "$url" -o "$tmpdir/$asset"
tar -xzf "$tmpdir/$asset" -C "$tmpdir" "$BINARY"
mkdir -p "$BINDIR"
install -m 0755 "$tmpdir/$BINARY" "$BINDIR/$BINARY"
echo "Installed $BINARY to $BINDIR/$BINARY"
