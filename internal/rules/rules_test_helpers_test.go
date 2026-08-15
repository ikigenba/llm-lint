package rules

import "fmt"

func testRuleSource(description, severity, include, exclude, prompt string) []byte {
	return []byte(fmt.Sprintf("---\ndescription: %s\nseverity: %s\ninclude: %s\nexclude: %s\n---\n%s", description, severity, include, exclude, prompt))
}
