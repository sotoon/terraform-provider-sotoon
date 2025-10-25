data "sotoon_iam_groups" "all" {
  workspace_id = "11111111-1111-1111-1111-111111111111"
}

output "all_group_names" {
  description = "A list of all group names in the workspace."
  value       = data.sotoon_iam_groups.all.groups.*.name
}