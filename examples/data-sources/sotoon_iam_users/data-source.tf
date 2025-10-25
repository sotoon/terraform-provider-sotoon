data "sotoon_iam_users" "all" {
  workspace_id = "11111111-1111-1111-1111-111111111111"
}

output "all_user_emails" {
  description = "A list of all user emails in the workspace."
  value       = data.sotoon_iam_users.all.users.*.email
}