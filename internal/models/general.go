package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type UUIDField struct {
	UUID types.String `tfsdk:"id"`
}

// A − B = {x ∈ A and x ∉ B}
func GetAdditionsOfUUIDAttributes(a []UUIDField, b []UUIDField) []UUIDField {
	diff := []UUIDField{}
	m := make(map[UUIDField]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}
