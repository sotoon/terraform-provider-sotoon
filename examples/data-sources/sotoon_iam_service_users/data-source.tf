data "sotoon_iam_service_users" "all" {
  # This data source retrieves all service users in the workspace
  # No additional parameters needed as it returns all service users
}

output "all_service_users" {
  value = data.sotoon_iam_service_users.all.users
}
