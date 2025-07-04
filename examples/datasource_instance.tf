
# Fetch data about all instances in the workspace
data "sotoon_instance" "all" {}

# You can now use this data, for example, to output the names of all running instances
output "running_instance_names" {
  description = "A list of all instance names that are currently in a running state."
  value = [
    for instance in data.sotoon_instance.all.instances : instance.name
    if instance.running == true
  ]
}
