
# Fetch data about all instances in the workspace
data "sotoon_instances" "all" {}

# You can now use this data, for example, to output the names of all running instances
output "running_instance_names" {
  description = "A list of all instance names that are currently in a running state."
  value = [
    for instance in data.sotoon_instances.all.instances : instance.name
    if instance.running == true
  ]
}

# Or find a specific instance and get its UID
locals {
  test_instance = one([
    for inst in data.sotoon_instances.all.instances : instance.name
    if inst.name == "test-1"
  ])
}

output "test_1_instance_uid" {
  description = "The UID of the instance named 'test-1'."
  value       = local.test_instance.uid
}