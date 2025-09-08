data "sotoon_iam_groups" "all" {
  # This data source retrieves all groups in the workspace
  # No additional parameters needed as it returns all groups
}

output "all_groups" {
  value = data.sotoon_iam_groups.all.groups
}
