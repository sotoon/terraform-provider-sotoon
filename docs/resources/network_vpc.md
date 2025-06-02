---
page_title: "sotoon_vpc Resource - terraform-provider-sotoon"
subcategory: "Network"
description: |-
  Creates and manages a Virtual Private Cloud (VPC) for your private networking infrastructure.
---

# sotoon_vpc (Resource)

Creates and manages a Virtual Private Cloud (VPC) for your private networking infrastructure.

## Example Usage

```terraform
resource "sotoon_vpc" "main" {
  name        = "production-vpc"
  cidr_block  = "10.0.0.0/16"
  region      = "neda"
  
  tags = {
    Environment = "production"
    Team        = "infrastructure"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (String, Required) - The name of the VPC.
* `cidr_block` - (String, Required) - The CIDR block for the VPC (e.g., "10.0.0.0/16").
* `region` - (String, Optional) - The region where the VPC will be created. Defaults to provider region.
* `enable_dns_hostnames` - (Boolean, Optional) - Enable DNS hostnames in the VPC. Defaults to true.
* `enable_dns_resolution` - (Boolean, Optional) - Enable DNS resolution in the VPC. Defaults to true.
* `tags` - (Map, Optional) - A map of tags to assign to the resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the VPC.
* `arn` - The ARN of the VPC.
* `default_security_group_id` - The ID of the default security group.
* `created_at` - The creation timestamp.
