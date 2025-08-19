data "sotoon_iam_service_user_groups" "developer" {
  workspace_id = var.workspace_id_target
  group_id     = local.target_group_developer.id
}

output "all_group_service_users_list" {
  description = "A list of all service-user names in the developer group."
  value       = data.sotoon_iam_service_user_groups.developer.service_users.*.name
}
