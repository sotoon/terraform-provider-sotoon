package provider

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnituniqueSortedEmptyString(t *testing.T) {
	in := make([]string, 0)
	if got := uniqueSorted(in); len(got) != 0 {
		t.Fatalf("uniqueSorted for empty input expect to return empty slice but returned %q", got)
	}
}

func TestUnituniqueSortedRemoveDuplicatesStrings(t *testing.T) {
	in := []string{"a", "b", "a"}
	expect := []string{"a", "b"}
	if got := uniqueSorted(in); !reflect.DeepEqual(got, expect) {
		t.Fatalf("uniqueSorted expect return %q but returned %q", expect, got)
	}
}

func TestUnituniqueSortedSortInput(t *testing.T) {
	in := []string{"b", "a"}
	expect := []string{"a", "b"}
	if got := uniqueSorted(in); !reflect.DeepEqual(got, expect) {
		t.Fatalf("uniqueSorted expect return %q but returned %q", expect, got)
	}
}

func TestUnithashOfIDs(t *testing.T) {
	in := []string{"1", "2"}
	expect := "17f8af97ad4a7f7639a4c9171d5185cbafb85462877a4746c21bdb0a4f940ca0"
	if got := hashOfIDs(in); got != expect {
		t.Fatalf("hashOfIDs expect return %q but returned %q", expect, got)
	}
}

func TestUnitfromSchemaSetToStringsNilInput(t *testing.T) {
	if got := fromSchemaSetToStrings(nil); got != nil {
		t.Fatal("fromSchemaSetToStrings expect return nil on nil input")
	}
}

func TestUnitfromSchemaSetToStrings(t *testing.T) {
	in := schema.NewSet(schema.HashString, []interface{}{"a", "b", "c", "a"})
	got := fromSchemaSetToStrings(in)
	for _, v := range got {
		if !in.Contains(v) {
			t.Fatalf("fromSchemaSetToStrings expect contains %q but not(%q)", v, got)
		}
	}
}

func TestUnitHashOfBindingIncludesItems(t *testing.T) {
	ids := []string{"user-b", "user-a"}
	first := hashOfBinding(ids, map[string]string{"namespace": "first", "cluster": "main"})
	same := hashOfBinding([]string{"user-a", "user-b"}, map[string]string{"cluster": "main", "namespace": "first"})
	second := hashOfBinding(ids, map[string]string{"namespace": "second", "cluster": "main"})
	if first != same {
		t.Fatalf("equivalent bindings produced different hashes: %q != %q", first, same)
	}
	if first == second {
		t.Fatalf("different item sets produced the same hash: %q", first)
	}
}

func TestUnitRemoveRoleBindingItemsPreservesOtherSets(t *testing.T) {
	first := map[string]string{"namespace": "first"}
	second := map[string]string{"namespace": "second"}
	remaining, found := removeRoleBindingItems(true, []map[string]string{first, second}, second)
	if !found {
		t.Fatal("wanted item set was not found")
	}
	if len(remaining) != 1 || !roleItemsEqual(remaining[0], first) {
		t.Fatalf("remaining item sets = %#v, want only %#v", remaining, first)
	}
}

func TestUnitRoleBindingHasEmptyItems(t *testing.T) {
	if !roleBindingHasItems(true, nil, map[string]string{}) {
		t.Fatal("an itemless remote binding did not match empty Terraform items")
	}
	if roleBindingHasItems(false, nil, map[string]string{}) {
		t.Fatal("an absent binding matched empty Terraform items")
	}
}
