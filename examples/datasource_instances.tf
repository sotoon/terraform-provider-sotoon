data "sotoon_instances" "all" {}

output "running_instance_names" {
  description = "A list of all instance names that are currently in a running state."
  value = [
    for instance in data.sotoon_instances.all.instances  : instance.name 
  ]
}

