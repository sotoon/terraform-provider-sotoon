data "sotoon_iam_service_user_roles" "all" {
  workspace_id    = "11111111-1111-1111-1111-111111111111"
  service_user_id = "44444444-4444-4444-4444-444444444444"
}

output "all_service_user_role_names" {
  description = "A list of all role names bound to the service user."
  value       = data.sotoon_iam_service_user_roles.all.roles.*.name
}
