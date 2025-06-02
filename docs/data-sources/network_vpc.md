---
page_title: "sotoon_vpc Data Source - terraform-provider-sotoon"
subcategory: "Network"
description: |-
  Get information about a Sotoon VPC.
---

# sotoon_vpc (Data Source)

Use this data source to get information about a Sotoon VPC for use in other resources.

## Example Usage

### Get VPC by Name

```terraform
data "sotoon_vpc" "main" {
  name = "production-vpc"
}

resource "sotoon_subnet" "example" {
  name              = "example-subnet"
  vpc_id            = data.sotoon_vpc.main.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "us-west-1a"
}
```

### Get VPC by ID

```terraform
data "sotoon_vpc" "main" {
  id = "vpc-12345678"
}
```

### Get VPC by Tags

```terraform
data "sotoon_vpc" "production" {
  tags = {
    Environment = "production"
    Team        = "infrastructure"
  }
}
```

## Argument Reference

The following arguments are supported:

* `id` (String, Optional) - The ID of the VPC to retrieve.
* `name` (String, Optional) - The name of the VPC to retrieve.
* `tags` (Map, Optional) - A map of tags assigned to the VPC.

~> **NOTE:** Either `id` or `name` must be specified. If both are provided, `id` takes precedence.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the VPC.
* `cidr_block` - The CIDR block of the VPC.
* `region` - The region where the VPC is located.
* `enable_dns_hostnames` - Whether DNS hostnames are enabled in the VPC.
* `enable_dns_resolution` - Whether DNS resolution is enabled in the VPC.
* `default_security_group_id` - The ID of the default security group.
* `created_at` - The creation timestamp of the VPC.
* `tags` - A map of tags assigned to the VPC.
```