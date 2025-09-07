data "sotoon_iam_roles" "all" {
  # This data source retrieves all roles in the workspace
  # No additional parameters needed as it returns all roles
}

output "all_roles" {
  value = data.sotoon_iam_roles.all.roles
}
