package rules

import (
	"errors"
	"strings"
	"testing"
)

// R-2NWS-3P5X
func TestSelectUnknownDisabledIDReturnsNamedRuleError(t *testing.T) {
	_, err := Select([]Rule{{ID: "known"}}, nil, []string{"missing-rule"})
	if !errors.Is(err, ErrRule) || !strings.Contains(err.Error(), "missing-rule") {
		t.Fatalf("Select() error = %v, want ErrRule naming unknown disabled id", err)
	}
}
