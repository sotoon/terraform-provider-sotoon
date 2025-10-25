data "sotoon_iam_group_roles" "developer" {
  workspace_id = "11111111-1111-1111-1111-111111111111"
  group_id     = "33333333-3333-3333-3333-333333333333"
}

output "all_group_roles_list" {
  description = "A list of all role names in the developer group."
  value       = data.sotoon_iam_group_roles.developer.roles.*.name
}
