package rules

import (
	"reflect"
	"testing"
)

// R-GD2Y-AJCT
func TestSelectReturnsEnabledSubsetInEnableOrder(t *testing.T) {
	first := Rule{ID: "first"}
	second := Rule{ID: "second"}
	got, err := Select([]Rule{first, second}, []string{"second"})
	if err != nil {
		t.Fatalf("Select() error = %v", err)
	}
	if !reflect.DeepEqual(got, []Rule{second}) {
		t.Fatalf("Select() = %#v, want only %#v", got, second)
	}

	empty, err := Select([]Rule{first, second}, nil)
	if err != nil || len(empty) != 0 {
		t.Fatalf("Select(nil enable) = %#v, %v; want empty, nil", empty, err)
	}
}
