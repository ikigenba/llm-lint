package rules

import (
	"errors"
	"strings"
	"testing"
)

// R-GGQN-FUKW
func TestSelectUnknownEnabledIDReturnsNamedRuleError(t *testing.T) {
	_, err := Select([]Rule{{ID: "known"}}, []string{"missing-rule"}, nil)
	if !errors.Is(err, ErrRule) || !strings.Contains(err.Error(), "missing-rule") {
		t.Fatalf("Select() error = %v, want ErrRule naming unknown id", err)
	}
}
