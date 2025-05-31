# sotoon_virtual_machine

Get information about a Sotoon Virtual Machine.

## Example Usage

```
data "sotoon_virtual_machine" "example" {
  hostname = "web-server-01"
}

# Or find by ID
data "sotoon_virtual_machine" "by_id" {
  id = "vm-1234567890abcdef"
}

# Use the data source
resource "sotoon_firewall_rule" "example" {
  firewall_id = sotoon_firewall.web.id
  protocol    = "tcp"
  port        = 80
  source_ip   = data.sotoon_virtual_machine.example.private_ip
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The unique identifier of the virtual machine.

* `hostname` - (Optional) The hostname of the virtual machine.

~> **Note:** Either `id` or `hostname` must be specified, but not both.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `zone` - The availability zone where the virtual machine is deployed.

* `am_id` - The allocation model ID defining the size/resources of the virtual machine.

* `os` - The operating system image deployed on the virtual machine.

* `backup` - Whether automated backups are enabled for this virtual machine.

* `vpc_id` - The ID of the VPC where the virtual machine is deployed.

* `ssh_keys` - List of SSH public keys added to the virtual machine.

* `user_data` - Cloud-init user data used for virtual machine customization.

* `tags` - Map of tags assigned to the virtual machine.

* `firewall_ids` - List of firewall IDs attached to this virtual machine.

* `volume_size` - Size of the root volume in GB.

* `disk_type` - Type of disk used for the root volume.

* `private_ip` - The primary private IP address of the virtual machine.

* `public_ip` - The public IP address of the virtual machine (if enabled).

* `status` - The current status of the virtual machine.

* `created_at` - The timestamp when the virtual machine was created.
