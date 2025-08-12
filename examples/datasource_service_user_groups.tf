data "sotoon_iam_service_user_groups" "all" {
  workspace_id = var.workspace_id_target
  group_id     = "3f355c0a-7871-47e4-8dcc-2d752f55b720"
}

output "all_group_service_users_list" {
  description = "A list of all service-user names in the group."
  value       = data.sotoon_iam_service_user_groups.all.service_users.*.name
}