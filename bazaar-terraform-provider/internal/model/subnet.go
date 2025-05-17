package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type Subnets struct {
	Subnets []Subnet `tfsdk:"subnets"`
}

type Subnet struct {
	Name      types.String  `tfsdk:"name"`
	GatewayIP types.String  `tfsdk:"gateway_ip"`
	Cidr      types.String  `tfsdk:"cidr"`
	Routes    []SubnetRoute `tfsdk:"routes"`
}

type SubnetRoute struct {
	To         types.String `tfsdk:"to"`
	ExternalIP types.String `tfsdk:"external_ip"`
}
