data "sotoon_iam_user_roles" "target" {
  user_id = local.target_user.id
}

output "user_role_names" {
  description = "All role names bound to the user."
  value       = data.sotoon_iam_user_roles.target.roles.*.name
}
