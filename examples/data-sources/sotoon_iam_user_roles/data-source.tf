data "sotoon_iam_user_roles" "target" {
  user_id = "22222222-2222-2222-2222-222222222222"
}

output "user_role_names" {
  description = "All role names bound to the user."
  value       = data.sotoon_iam_user_roles.target.roles.*.name
}

