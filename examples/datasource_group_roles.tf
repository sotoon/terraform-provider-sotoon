data "sotoon_iam_group_roles" "all" {
  workspace_id = var.workspace_id_target
  group_id     = local.target_group_developer.id
}

output "all_group_roles_names" {
  description = "A list of all role names in the group."
  value       = data.sotoon_iam_group_roles.all.roles.*.name
}
