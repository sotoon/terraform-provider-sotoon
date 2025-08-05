data "sotoon_iam_groups" "all" {
  workspace_id = var.workspace_id_target
}

output "all_group_names" {
  description = "A list of all group names in the workspace."
  value       = data.sotoon_iam_groups.all.groups.*.name
}