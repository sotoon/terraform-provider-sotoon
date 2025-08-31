data "sotoon_iam_group_roles" "developer" {
  workspace_id = var.workspace_id_target
  group_id     = local.target_group_developer.id
}

output "all_group_roles_list" {
  description = "A list of all role names in the developer group."
  value       = data.sotoon_iam_group_roles.developer.roles.*.name
}
