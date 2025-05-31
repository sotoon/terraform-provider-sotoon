# sotoon_allocation_models

Get information about available Sotoon allocation models (VM sizes).

## Example Usage

```
data "sotoon_allocation_models" "all" {}

# Filter by type
data "sotoon_allocation_models" "standard" {
  type = "standard"
}

# Output available models
output "available_models" {
  value = data.sotoon_allocation_models.all.models
}

# Use in a virtual machine resource
resource "sotoon_virtual_machine" "example" {
  hostname = "test-vm"
  zone     = "neda"
  am_id    = data.sotoon_allocation_models.standard.models[0].id
  os       = "ubuntu-22.04"
  vpc_id   = sotoon_vpc.main.id
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Optional) Filter allocation models by type. Available values: `standard`, `performance`, `memory`.

* `min_vcpus` - (Optional) Filter allocation models with at least this many vCPUs.

* `min_memory` - (Optional) Filter allocation models with at least this much memory (in GB).

* `max_vcpus` - (Optional) Filter allocation models with at most this many vCPUs.

* `max_memory` - (Optional) Filter allocation models with at most this much memory (in GB).

## Attribute Reference

The following attributes are exported:

* `models` - List of allocation models matching the filter criteria. Each model contains:
* `id` - The allocation model ID (e.g., `standard-2x4`).
* `vcpus` - Number of virtual CPUs.
* `memory` - Amount of memory in GB.
* `type` - Type of the allocation model (`standard`, `performance`, `memory`).
* `description` - Human-readable description of the allocation model.
