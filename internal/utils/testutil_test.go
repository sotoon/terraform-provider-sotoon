package testutil

import (
	"regexp"
	"testing"
)

func TestUnitBaseProviderBlock(t *testing.T) {
	expected := `
provider "sotoon" {}
`
	if got := BaseProviderBlock(); got != expected{
		t.Fatalf("BaseProviderBlock() = %q, want %q", got, expected)
	}
}

func TestUnitRandEmail(t *testing.T) {
	pattern := `^tf-acc-[a-zA-Z0-9]+@example\.test$`
	re := regexp.MustCompile(pattern)

	if got := RandEmail(); !re.MatchString(got){
		t.Fatalf("RandEmail() = %q, pattern must be like: %q", got, pattern);
	}
}

func TestUnitRandName(t *testing.T) {
	pattern := `^random-[a-zA-Z0-9]{8}$`
	re := regexp.MustCompile(pattern)

	if got := RandName("random"); !re.MatchString(got){
		t.Fatalf("RandName('random') = %q, pattern must be like: %q", got, pattern);
	}  
}
