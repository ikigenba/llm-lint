package rules

import (
	"reflect"
	"testing"
)

// R-4HI3-PKJF
func TestParseOmittedIncludeUsesDefaultCodeGlobs(t *testing.T) {
	src := []byte("---\ndescription: general code check\nseverity: error\nexclude: []\n---\nbody")
	rule, err := Parse("general-code.md", src)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	want := []string{
		"**/*.go", "**/*.py", "**/*.js", "**/*.ts", "**/*.java", "**/*.rb",
		"**/*.rs", "**/*.c", "**/*.h", "**/*.cpp", "**/*.cs",
	}
	if !reflect.DeepEqual(DefaultCodeGlobs, want) {
		t.Fatalf("DefaultCodeGlobs = %v, want %v", DefaultCodeGlobs, want)
	}
	if !reflect.DeepEqual(rule.Include, DefaultCodeGlobs) {
		t.Fatalf("Parse() Include = %v, want DefaultCodeGlobs %v", rule.Include, DefaultCodeGlobs)
	}
}
