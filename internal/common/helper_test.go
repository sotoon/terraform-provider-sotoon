package common

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"testing"
)

func TestUnitUniqueSortedEmptyString(t *testing.T) {
	in := make([]string, 0)

	if got := UniqueSorted(in); len(got) != 0 {
		t.Fatalf("UniqueSorted for empty input expect to return empty slice but returned %q", got)
	}
}

func TestUnitUniqueSortedRemoveDuplicatesStrings(t *testing.T) {
	in := []string{"a", "b", "a"}
	expect := []string{"a", "b"}

	if got := UniqueSorted(in); reflect.DeepEqual(got, expect) {
		t.Fatalf("UniqueSorted expect return %q but returned %q", expect, got)
	}
}

func TestUnitUniqueSortedSortInput(t *testing.T) {
	in := []string{"b", "a"}
	expect := []string{"a", "b"}

	if got := UniqueSorted(in); !reflect.DeepEqual(got, expect) {
		t.Fatalf("UniqueSorted expect return %q but returned %q", expect, got)
	}
}

func TestUnitHashOfIDs(t *testing.T) {
	in := []string{"1", "2"}
	expect := "17f8af97ad4a7f7639a4c9171d5185cbafb85462877a4746c21bdb0a4f940ca0" // hash of "1,2" in byte
	if got := HashOfIDs(in); got != expect {
		t.Fatalf("HashOfIDs expect return %q but returned %q", expect, got)
	}
}

func TestUnitFromSchemaSetToStringsNilInput(t *testing.T) {

	if got := FromSchemaSetToStrings(nil); got != nil{
		t.Fatalf("FromSchemaSetToStrings expect return nil on nil input")
	}
}

func TestUnitFromSchemaSetToStrings(t *testing.T) {

	in := schema.NewSet(schema.HashString, []interface{}{"a", "b", "c", "a"})
	got := FromSchemaSetToStrings(in);

	for _, v := range got{
		if !in.Contains(v){
			t.Fatalf("FromSchemaSetToStrings expect contains %q but not(%q)", v, got)
		} 
	} 
}