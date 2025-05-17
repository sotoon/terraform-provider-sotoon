package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type Computes struct {
	Computes []Compute `tfsdk:"computes"`
}

type Compute struct {
	Name       types.String `tfsdk:"name"`
	IAMEnabled types.Bool   `tfsdk:"iam_enabled"`
	Image      types.String `tfsdk:"image"`
	Username   types.String `tfsdk:"username"`
	Size       types.String `tfsdk:"size"`
	Volumes    []Volume     `tfsdk:"volumes"`
	PoweredOn  types.Bool   `tfsdk:"powered_on"`
	Subnet     types.String `tfsdk:"subnet"`
	ExternalIP types.String `tfsdk:"external_ip"`
}

type Volume struct {
	LocalDisk  LocalDisk  `tfsdk:"local_disk"`
	Pvc        Pvc        `tfsdk:"pvc"`
	RemoteDisk RemoteDisk `tfsdk:"remote_disk"`
}

type LocalDisk struct {
	DiskSize types.Int64  `tfsdk:"disk_size_gb"`
	Name     types.String `tfsdk:"name"`
}

type Pvc struct {
	Name types.String `tfsdk:"name"`
}

type RemoteDisk struct {
	Name types.String `tfsdk:"name"`
	Size types.Int64  `tfsdk:"size"`
	Tier types.String `tfsdk:"tier"`
}
