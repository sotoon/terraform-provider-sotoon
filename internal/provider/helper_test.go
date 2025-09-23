package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"testing"
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
	expect := "17f8af97ad4a7f7639a4c9171d5185cbafb85462877a4746c21bdb0a4f940ca0" // hash of "1,2" in byte
	if got := hashOfIDs(in); got != expect {
		t.Fatalf("hashOfIDs expect return %q but returned %q", expect, got)
	}
}

func TestUnitfromSchemaSetToStringsNilInput(t *testing.T) {

	if got := fromSchemaSetToStrings(nil); got != nil{
		t.Fatalf("fromSchemaSetToStrings expect return nil on nil input")
	}
}

func TestUnitfromSchemaSetToStrings(t *testing.T) {

	in := schema.NewSet(schema.HashString, []interface{}{"a", "b", "c", "a"})
	got := fromSchemaSetToStrings(in);

	for _, v := range got{
		if !in.Contains(v){
			t.Fatalf("fromSchemaSetToStrings expect contains %q but not(%q)", v, got)
		} 
	} 
}