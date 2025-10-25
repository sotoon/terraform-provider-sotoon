resource "sotoon_iam_group" "developer" {
  name = "developer"
  description = "description about this group"
}

output "group_id" {
  description = "The name of the created group of developer."
  value       = sotoon_iam_group.developer.name
}
