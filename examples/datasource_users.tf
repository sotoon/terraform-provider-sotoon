data "sotoon_iam_users" "all" {
  workspace_id = var.workspace_id_target
}

output "all_user_emails" {
  description = "A list of all user emails in the workspace."
  value       = data.sotoon_iam_users.all.users.*.email
}

output "moein_user_id" {
  description = "The UUID for the user with the specified email."
  value = one([
    for user in data.sotoon_iam_users.all.users : user.id if user.email == "moein.tavakoli@sotoon.ir"
  ])
}

output "first_user_details" {
  description = "All details for the first user in the list."
  value = one(data.sotoon_iam_users.all.users[*])
}
