data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

data "sotoon_iam_group" "deployers" {
  workspace_id = data.sotoon_workspace.mycompany.id
  name         = "deployers"
}

# Retrieve all service-users of the workspace
data "sotoon_iam_service_users" "all" {
  workspace_id = data.sotoon_workspace.mycompany.id
}

output "all_service_users" {
  description = "All service-users in the workspace"
  value       = data.sotoon_iam_service_users.all.service_users
}

# Retrieve all service-users of the "deployers" group
data "sotoon_iam_service_users" "deployers" {
  workspace_id = data.sotoon_workspace.mycompany.id
  group_id     = data.sotoon_iam_group.deployers.id
}

output "group_service_users" {
  description = "All service-users in the group"
  value       = data.sotoon_iam_service_users.deployers.service_users
}
