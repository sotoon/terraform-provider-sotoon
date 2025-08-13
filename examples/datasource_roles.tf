data "sotoon_iam_roles" "all" {
  workspace_id = var.workspace_id_target
}

output "all_role_names" {
  description = "A list of all role names in the workspace."
  value       = data.sotoon_iam_roles.all.roles.*.name
}