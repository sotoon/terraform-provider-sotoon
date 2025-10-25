data "sotoon_iam_service_user_groups" "developer" {
  workspace_id = "11111111-1111-1111-1111-111111111111"
  group_id     = "33333333-3333-3333-3333-333333333333"
}

output "all_group_service_users_list" {
  description = "A list of all service-user names in the developer group."
  value       = data.sotoon_iam_service_user_groups.developer.service_users.*.name
}
