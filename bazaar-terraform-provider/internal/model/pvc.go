package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type PVC struct {
	Name        types.String `tfsdk:"name"`
	Size        types.String `tfsdk:"size"`
	StorageTier types.String `tfsdk:"tier"`
}

type PVCS struct {
	PVCS []PVC `tfsdk:"pvcs"`
}
