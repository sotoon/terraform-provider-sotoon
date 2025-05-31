# sotoon_virtual_machine

Provides a Sotoon Virtual Machine resource. This resource allows you to create and manage virtual machines on the Sotoon cloud platform.

## Example Usage

```
resource "sotoon_virtual_machine" "web_server" {
  hostname     = "web-server-01"
  zone         = "neda"
  am_id        = "standard-2x4"
  os           = "ubuntu-22.04"
  backup       = true
  vpc_id       = sotoon_vpc.main.id
  disk_type    = "premium"
  volume_size  = 80
  
  ssh_keys = [
    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."
  ]
  
  user_data = <<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    systemctl start nginx
    systemctl enable nginx
  EOF
  
  firewall_ids = [
    sotoon_firewall.web.id,
    sotoon_firewall.ssh.id
  ]
  
  tags = {
    environment = "production"
    service     = "web"
    team        = "infrastructure"
  }
}
```

Provides a Sotoon Virtual Machine resource. This resource allows you to create and manage virtual machines on the Sotoon cloud platform.

## Example Usage

```
resource "sotoon_virtual_machine" "web_server" {
  hostname     = "web-server-01"
  zone         = "neda"
  am_id        = "standard-2x4"
  os           = "ubuntu-22.04"
  backup       = true
  vpc_id       = sotoon_vpc.main.id
  disk_type    = "premium"
  volume_size  = 80
  
  ssh_keys = [
    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."
  ]
  
  user_data = <<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    systemctl start nginx
    systemctl enable nginx
  EOF
  
  firewall_ids = [
    sotoon_firewall.web.id,
    sotoon_firewall.ssh.id
  ]
  
  tags = {
    environment = "production"
    service     = "web"
    team        = "infrastructure"
  }
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Required) The hostname of the virtual machine. Must be unique within the zone.
* `zone` - (Required) Availability zone where the virtual machine will be deployed. Available values: neda, afra. See Availability Zone Reference.
* `am_id` - (Required) The allocation model ID which defines the size/resources of the virtual machine. See Allocation Model Reference.
* `os` - (Required) Operating system image to deploy. See Operating System Reference.
* `vpc_id` - (Required) The ID of the VPC where the virtual machine will be deployed.
* `backup` - (Optional) Whether to enable automated backups for this virtual machine. Defaults to false.
* `ssh_keys` - (Optional) A list of SSH public keys to be added to the virtual machine for authentication.
* `user_data` - (Optional) Cloud-init user data to be used for customizing the virtual machine during creation. Must be base64 encoded or plain text.
* `tags` - (Optional) A map of tags to assign to the virtual machine for organization and billing purposes.
* `firewall_ids` - (Optional) List of firewall IDs to attach to this virtual machine for network security.
* `volume_size` - (Optional) Size of the root volume in GB. Defaults to the OS image default size. Minimum 20GB, maximum 2TB.
* `disk_type` - (Optional) Type of disk to use for the root volume. Available options: premium, standard, ultra, local. Defaults to standard.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the virtual machine.
* `private_ip` - The primary private IP address of the virtual machine within the VPC.
* `public_ip` - The public IP address of the virtual machine if enabled and available.
* `status` - The current status of the virtual machine (running, stopped, pending, terminated).
* `created_at` - The timestamp when the virtual machine was created (RFC3339 format).

## Allocation Model (am_id) Reference

| AM ID | vCPUs | RAM (GB) | Description |
|-------|-------|----------|-------------|
| `standard-1x2` | 1 | 2 | Standard VM with 1 vCPU and 2GB RAM |
| `standard-2x4` | 2 | 4 | Standard VM with 2 vCPUs and 4GB RAM |
| `standard-4x8` | 4 | 8 | Standard VM with 4 vCPUs and 8GB RAM |
| `standard-8x16` | 8 | 16 | Standard VM with 8 vCPUs and 16GB RAM |
| `performance-2x8` | 2 | 8 | Performance VM with 2 vCPUs and 8GB RAM |
| `performance-4x16` | 4 | 16 | Performance VM with 4 vCPUs and 16GB RAM |
| `memory-2x16` | 2 | 16 | Memory-optimized VM with 2 vCPUs and 16GB RAM |
| `memory-4x32` | 4 | 32 | Memory-optimized VM with 4 vCPUs and 32GB RAM |

## Operating System (OS) Reference

| OS ID | Description |
|-------|-------------|
| `ubuntu-20.04` | Ubuntu Server 20.04 LTS |
| `ubuntu-22.04` | Ubuntu Server 22.04 LTS |
| `debian-11` | Debian 11 (Bullseye) |
| `debian-12` | Debian 12 (Bookworm) |

## Availability Zone Reference

| Zone | Location |
|------|----------|
| `neda` | Primary data center |
| `afra` | Secondary data center |

## Disk Type Reference

| Disk Type | Description | Use Case |
|-----------|-------------|----------|
| `premium` | Premium SSD | High-performance workloads requiring consistent IOPS |
| `standard` | Standard SSD | General purpose workloads with balanced performance and cost |
| `ultra` | Ultra SSD | Mission-critical workloads requiring extremely high IOPS and throughput |
| `local` | Local SSD | Temporary storage with high performance for cache/scratch data |

* `create` - (Default `10m`) How long to wait for a virtual machine to be created.
* `update` - (Default `5m`) How long to wait for a virtual machine to be updated.
* `delete` - (Default `10m`) How long to wait for a virtual machine to be deleted.

## Error Handling

Common error codes and troubleshooting steps:

| Error Code | Description | Resolution |
|------------|-------------|------------|
| `auth_error` | Authentication failed | Verify your API token is correct and has proper permissions |
| `quota_exceeded` | Resource quota exceeded | Request a quota increase or delete unused resources |
| `invalid_am_id` | Invalid allocation model ID | Check the AM ID reference for valid values |
| `zone_unavailable` | Selected zone is unavailable | Try another zone or contact support |
| `resource_conflict` | Resource naming conflict | Use a different hostname or check for existing resources |
