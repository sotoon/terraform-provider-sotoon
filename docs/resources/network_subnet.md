---
page_title: "sotoon_subnet Resource - terraform-provider-sotoon"
subcategory: "Network"
description: |-
  Creates and manages subnets within a VPC for network segmentation.
---

# sotoon_subnet (Resource)

Creates and manages subnets within a VPC for network segmentation.

## Example Usage

### Private Database Subnet

```terraform
resource "sotoon_subnet" "database" {
  name               = "database-subnet"
  vpc_id             = sotoon_vpc.main.id
  cidr_block         = "10.0.1.0/24"
  availability_zone  = "neda"
  subnet_type        = "private"
  
  tags = {
    Tier = "database"
    Type = "private"
  }
}
```

### Private Subnet for Kubernetes

```
resource "sotoon_subnet" "kubernetes" {
  name               = "k8s-subnet"
  vpc_id             = sotoon_vpc.main.id
  cidr_block         = "10.0.2.0/24"
  availability_zone  = "neda"
  subnet_type        = "private"
  
  tags = {
    Tier = "application"
    Type = "kubernetes"
  }
}
```

### Public Subnet

```
resource "sotoon_subnet" "public" {
  name               = "public-subnet"
  vpc_id             = sotoon_vpc.main.id
  cidr_block         = "10.0.10.0/24"
  availability_zone  = "neda"
  subnet_type        = "public"
  auto_assign_public_ip = true
  
  tags = {
    Tier = "public"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the subnet. Must be unique within the VPC.
* `vpc_id` - (Required) The ID of the VPC where the subnet will be created.
* `cidr_block` - (Required) The CIDR block for the subnet (e.g., "10.0.1.0/24"). Must be a subset of the VPC's CIDR block and not overlap with existing subnets.
* `availability_zone` - (Required) The availability zone for the subnet. Must be available in the VPC's region.
* `subnet_type` - (Optional) The type of subnet: `"public"` or `"private"`. Defaults to `"private"`.
  * `"public"` - Subnet with internet gateway access
  * `"private"` - Subnet without direct internet access
* `auto_assign_public_ip` - (Optional) Assign public IP automatically for instances launched in this subnet. Defaults to `true` for public subnets, `false` for private subnets.
* `tags` - (Optional) A map of tags to assign to the resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the subnet.
* `arn` - The Amazon Resource Name (ARN) of the subnet.
* `available_ip_count` - The number of available IP addresses in the subnet.
* `created_at` - The timestamp when the subnet was created (RFC3339 format).
* `state` - The current state of the subnet (e.g., "available", "pending").
* `route_table_id` - The ID of the route table associated with the subnet.

## Timeouts

`sotoon_subnet` provides the following timeouts:

* `create` - (Default 5 minutes) Used when creating the subnet
* `update` - (Default 3 minutes) Used when updating the subnet
* `delete` - (Default 5 minutes) Used when deleting the subnet
