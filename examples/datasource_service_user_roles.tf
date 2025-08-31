data "sotoon_iam_service_user_roles" "all" {
  workspace_id    = var.workspace_id_target
  service_user_id = local.target_service_user_builder.id
}

output "all_service_user_role_names" {
  description = "A list of all role names bound to the service user."
  value       = data.sotoon_iam_service_user_roles.all.roles.*.name
}
