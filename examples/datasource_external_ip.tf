data "sotoon_external_ips" "all" {}

output "all_ips_data" {
  description = "The complete list of data for all external IPs."
  value       = data.sotoon_external_ips.all.items
}


# Get only available (unbound) IPs

data "sotoon_external_ips" "available" {
  filter {
    unbound_only = true
  }
}

output "all_available_ips_data" {
  description = "The complete list of available external IPs."
  value       = data.sotoon_external_ips.available.items
}
