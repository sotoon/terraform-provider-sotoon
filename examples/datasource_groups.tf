data "sotoon_iam_groups" "all" {
  workspace_id = "4dc150c4-8e77-480f-98b6-821f8fc50227"
}

output "all_group_names" {
  description = "A list of all group names in the workspace."
  value       = data.sotoon_iam_groups.all.groups.*.name
}