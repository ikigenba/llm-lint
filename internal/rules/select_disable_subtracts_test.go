package rules

import (
	"reflect"
	"testing"
)

// R-2MOV-PXF8
func TestSelectDisableSubtractsFromDefaultAndEnabledRules(t *testing.T) {
	first := Rule{ID: "first"}
	second := Rule{ID: "second"}
	all := []Rule{first, second}

	tests := []struct {
		name   string
		enable []string
		want   []Rule
	}{
		{name: "all-on default", want: []Rule{first}},
		{name: "explicit allowlist", enable: []string{"second", "first"}, want: []Rule{first}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Select(all, test.enable, []string{"second"})
			if err != nil {
				t.Fatalf("Select() error = %v", err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("Select() = %#v, want %#v", got, test.want)
			}
		})
	}
}
