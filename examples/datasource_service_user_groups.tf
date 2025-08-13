data "sotoon_iam_service_user_groups" "all" {
  workspace_id = var.workspace_id_target
  group_id     = local.target_group_developer.id
}

output "all_groups_service_users_list" {
  description = "A list of all service-user names in the group."
  value       = data.sotoon_iam_service_user_groups.all.service_users.*.name
}