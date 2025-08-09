data "sotoon_iam_group_users_list" "developer" {
  workspace_id = var.workspace_id_target
  group_id     = local.target_group_developer.id
}

output "all_group_users_list" {
  description = "A list of all user names in the developer group."
  value       = data.sotoon_iam_group_users_list.developer.users.*.name
}
