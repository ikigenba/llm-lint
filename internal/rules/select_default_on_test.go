package rules

import (
	"reflect"
	"testing"
)

// R-2LGZ-C5OJ
func TestSelectDefaultsToEveryRuleInCatalogOrder(t *testing.T) {
	all := []Rule{{ID: "first"}, {ID: "second"}, {ID: "third"}}

	got, err := Select(all, nil, nil)
	if err != nil {
		t.Fatalf("Select() error = %v", err)
	}
	if !reflect.DeepEqual(got, all) {
		t.Fatalf("Select() = %#v, want every rule in catalog order %#v", got, all)
	}
}
