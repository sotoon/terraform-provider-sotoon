data "sotoon_iam_group_users" "all" {
  workspace_id = "11111111-1111-1111-1111-111111111111"
  group_id     = "33333333-3333-3333-3333-333333333333"
}

output "all_group_users_list" {
  description = "A list of all users names in the group."
  value       = data.sotoon_iam_group_users.all.users.*.name
}