package rules

import (
	"reflect"
	"testing"
)

// R-GD2Y-AJCT
func TestSelectReturnsEnabledSubsetInEnableOrder(t *testing.T) {
	first := Rule{ID: "first"}
	second := Rule{ID: "second"}
	got, err := Select([]Rule{first, second}, []string{"second"}, nil)
	if err != nil {
		t.Fatalf("Select() error = %v", err)
	}
	if !reflect.DeepEqual(got, []Rule{second}) {
		t.Fatalf("Select() = %#v, want only %#v", got, second)
	}
}
