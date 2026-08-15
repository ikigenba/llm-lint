package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func readReleaseFile(t *testing.T, root, name string) string {
	t.Helper()

	contents, err := os.ReadFile(filepath.Join(root, name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(contents)
}

func requireReleaseFragments(t *testing.T, name, contents string, fragments ...string) {
	t.Helper()

	for _, fragment := range fragments {
		if !strings.Contains(contents, fragment) {
			t.Errorf("%s does not contain %q", name, fragment)
		}
	}
}

// R-HJRF-OIB4
func TestReleaseFilesAgreeOnArtifactAndRepository(t *testing.T) {
	root := filepath.Join("..", "..")
	goreleaser := readReleaseFile(t, root, ".goreleaser.yaml")
	workflow := readReleaseFile(t, root, filepath.Join(".github", "workflows", "release.yml"))
	installer := readReleaseFile(t, root, "install.sh")

	requireReleaseFragments(t, ".goreleaser.yaml", goreleaser,
		"version: 2",
		"project_name: llm-lint",
		"main: ./cmd/llm-lint",
		"binary: llm-lint",
		"env: [CGO_ENABLED=0]",
		"goos: [linux, darwin]",
		"goarch: [amd64, arm64]",
		`ldflags: ["-s -w -X main.version={{ .Tag }}"]`,
		"formats: [tar.gz]",
		`name_template: "{{ .Binary }}_{{ .Os }}_{{ .Arch }}"`,
		`name_template: "checksums.txt"`,
		"owner: ikigenba",
		"name: llm-lint",
	)
	requireReleaseFragments(t, ".github/workflows/release.yml", workflow,
		"tags: [\"v*\"]",
		"contents: write",
		"fetch-depth: 0",
		"go-version: \"1.26\"",
		"goreleaser/goreleaser-action@v6",
		"args: \"release --clean\"",
	)
	requireReleaseFragments(t, "install.sh", installer,
		"REPO=ikigenba/llm-lint",
		"BINARY=llm-lint",
		"Linux) os=linux",
		"Darwin) os=darwin",
		"x86_64 | amd64) arch=amd64",
		"arm64 | aarch64) arch=arm64",
		`asset="${BINARY}_${os}_${arch}.tar.gz"`,
	)
}

// R-HKZC-2A1T
func TestInstallerSyntaxAndMakeInstallContract(t *testing.T) {
	root := filepath.Join("..", "..")
	installerPath := filepath.Join(root, "install.sh")
	command := exec.Command("sh", "-n", installerPath)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("sh -n install.sh: %v\n%s", err, output)
	}

	installer := readReleaseFile(t, root, "install.sh")
	requireReleaseFragments(t, "install.sh", installer,
		"LLM_LINT_VERSION=${LLM_LINT_VERSION:-latest}",
		"BINDIR=${BINDIR:-${PREFIX:-$HOME/.local}/bin}",
		"releases/latest/download/${asset}",
		"releases/download/${LLM_LINT_VERSION}/${asset}",
		`install -m 0755 "$tmpdir/$BINARY" "$BINDIR/$BINARY"`,
	)

	makefile := readReleaseFile(t, root, "Makefile")
	requireReleaseFragments(t, "Makefile", makefile,
		"build:",
		"go build",
		"-o bin/llm-lint ./cmd/llm-lint",
		"install: build",
		"install -m 0755 bin/llm-lint $(DESTDIR)$(PREFIX)/bin/llm-lint",
	)
}
