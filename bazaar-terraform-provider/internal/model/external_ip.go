package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type ExternalIPSpec struct {
	IP       string `json:"ip"`
	Reserved bool   `json:"reserved"`
	Gateway  string `json:"gateway"`
}

type ExternalIP struct {
	Name      types.String `tfsdk:"name"`
	GatewayIP types.String `tfsdk:"gateway_ip"`
	IP        types.String `tfsdk:"ip"`
	Reserved  types.Bool   `tfsdk:"reserved"`
}
