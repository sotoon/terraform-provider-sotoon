data "sotoon_iam_group_details" "developer" {
  workspace_id = var.workspace_id_target
  group_id     = local.target_group_developer.id
}

output "group_name" {
  description = "The name of the group"
  value       = local.target_group_developer.name
}
