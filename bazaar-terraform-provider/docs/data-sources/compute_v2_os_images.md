# sotoon_os_images

Get information about available Sotoon operating system images.

## Example Usage

```
data "sotoon_os_images" "all" {}

# Filter by OS family
data "sotoon_os_images" "ubuntu" {
  family = "ubuntu"
}

# Filter by architecture
data "sotoon_os_images" "x86_64" {
  architecture = "x86_64"
}

# Get latest Ubuntu image
data "sotoon_os_images" "latest_ubuntu" {
  family = "ubuntu"
  latest = true
}

# Output available images
output "available_images" {
  value = data.sotoon_os_images.all.images
}

# Use in a virtual machine resource
resource "sotoon_virtual_machine" "example" {
  hostname = "test-vm"
  zone     = "neda"
  am_id    = "standard-2x4"
  os       = data.sotoon_os_images.latest_ubuntu.images[0].id
  vpc_id   = sotoon_vpc.main.id
}
```

## Argument Reference

The following arguments are supported:

* `family` - (Optional) Filter OS images by family. Available values: `ubuntu`, `debian`.

* `version` - (Optional) Filter OS images by specific version (e.g., `22.04`, `11`).

* `architecture` - (Optional) Filter OS images by architecture. Defaults to `x86_64`.

* `latest` - (Optional) If true, only return the latest version of each OS family. Defaults to `false`.

## Attribute Reference

The following attributes are exported:

* `images` - List of OS images matching the filter criteria. Each image contains:
* `id` - The OS image ID (e.g., `ubuntu-22.04`).
* `name` - Human-readable name of the OS image.
* `family` - OS family (`ubuntu`, `debian`).
* `version` - OS version.
* `architecture` - System architecture.
* `description` - Detailed description of the OS image.
* `min_disk_size` - Minimum required disk size in GB.
* `created_at` - Timestamp when the image was created.
