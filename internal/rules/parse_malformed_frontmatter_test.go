package rules

import (
	"errors"
	"strings"
	"testing"
)

// R-GBV1-WRM4
func TestParseMalformedFrontmatterReturnsNamedRuleError(t *testing.T) {
	tests := map[string]string{
		"unterminated":    "---\ndescription: test\nseverity: error\ninclude: [\"**/*\"]\n",
		"unknown-key":     "---\ndescription: test\nseverity: error\ninclude: [\"**/*\"]\nsurprise: yes\n---\nbody",
		"bad-severity":    "---\ndescription: test\nseverity: fatal\ninclude: [\"**/*\"]\n---\nbody",
		"missing-include": "---\ndescription: test\nseverity: error\n---\nbody",
		"empty-include":   "---\ndescription: test\nseverity: error\ninclude: []\n---\nbody",
	}
	for name, src := range tests {
		t.Run(name, func(t *testing.T) {
			filename := name + ".md"
			_, err := Parse(filename, []byte(src))
			if !errors.Is(err, ErrRule) || !strings.Contains(err.Error(), filename) {
				t.Fatalf("Parse() error = %v, want ErrRule naming %q", err, filename)
			}
		})
	}
}

func TestParseRejectsInvalidFilenameID(t *testing.T) {
	_, err := Parse("Bad_Name.md", testRuleSource("test", "error", `["**/*"]`, `[]`, "body"))
	if !errors.Is(err, ErrRule) || !strings.Contains(err.Error(), "Bad_Name.md") {
		t.Fatalf("Parse() error = %v, want named ErrRule", err)
	}
}
