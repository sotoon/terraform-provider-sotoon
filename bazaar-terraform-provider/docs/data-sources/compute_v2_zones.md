# sotoon_zones

Get information about available Sotoon availability zones.

## Example Usage

```
data "sotoon_zones" "all" {}

# Filter by status
data "sotoon_zones" "available" {
  status = "available"
}

# Output available zones
output "available_zones" {
  value = data.sotoon_zones.all.zones
}

# Use in a virtual machine resource
resource "sotoon_virtual_machine" "example" {
  hostname = "test-vm"
  zone     = data.sotoon_zones.available.zones[0].id
  am_id    = "standard-2x4"
  os       = "ubuntu-22.04"
  vpc_id   = sotoon_vpc.main.id
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Optional) Filter zones by status. Available values: `available`, `maintenance`, `unavailable`.

## Attribute Reference

The following attributes are exported:

* `zones` - List of availability zones matching the filter criteria. Each zone contains:
* `id` - The zone identifier (e.g., `neda`, `afra`).
* `name` - Human-readable name of the zone.
* `location` - Geographic location description.
* `status` - Current status of the zone.
* `capabilities` - List of capabilities supported in this zone.
