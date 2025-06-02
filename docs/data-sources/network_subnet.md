---
page_title: "sotoon_subnet Data Source - terraform-provider-sotoon"
subcategory: "Network"
description: |-
  Get information about a Sotoon subnet.
---

# sotoon_subnet (Data Source)

Use this data source to get information about a Sotoon subnet for use in other resources.

## Example Usage

### Get Subnet by Name

```terraform
data "sotoon_subnet" "database" {
  name = "database-subnet"
}

resource "sotoon_instance" "db" {
  subnet_id = data.sotoon_subnet.database.id
  # ... other configuration
}
```

### Get Subnet by ID

```terraform
data "sotoon_subnet" "database" {
  id = "subnet-12345678"
}
```

### Get Subnet by VPC and Availability Zone

```terraform
data "sotoon_subnet" "az_specific" {
  vpc_id            = "vpc-12345678"
  availability_zone = "us-west-1a"
  subnet_type       = "private"
}
```

### Get Subnet by Tags

```terraform
data "sotoon_subnet" "kubernetes" {
  tags = {
    Type = "kubernetes"
    Tier = "application"
  }
}
```

## Argument Reference

The following arguments are supported:

* `id` (String, Optional) - The ID of the subnet to retrieve.
* `name` (String, Optional) - The name of the subnet to retrieve.
* `vpc_id` (String, Optional) - The ID of the VPC that contains the subnet.
* `availability_zone` (String, Optional) - The availability zone of the subnet.
* `subnet_type` (String, Optional) - The type of subnet ("public" or "private").
* `cidr_block` (String, Optional) - The CIDR block of the subnet.
* `tags` (Map, Optional) - A map of tags assigned to the subnet.

~> **NOTE:** At least one of `id`, `name`, or a combination of `vpc_id` with other filters must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the subnet.
* `available_ip_count` - The number of available IP addresses in the subnet.
* `auto_assign_public_ip` - Whether public IP addresses are automatically assigned.
* `created_at` - The creation timestamp of the subnet.
* `tags` - A map of tags assigned to the subnet.

