data "sotoon_iam_group_details" "developer" {
  workspace_id = "11111111-1111-1111-1111-111111111111"
  group_id     = "33333333-3333-3333-3333-333333333333"
}

output "group_name" {
  description = "The name of the group"
  value       = data.sotoon_iam_group_details.developer.name
}
