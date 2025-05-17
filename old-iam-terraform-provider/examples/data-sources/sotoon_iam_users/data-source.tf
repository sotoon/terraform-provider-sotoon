data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

data "sotoon_iam_group" "deployers" {
  workspace_id = data.sotoon_workspace.mycompany.id
  name         = "deployers"
}

# Retrieve all users of the workspace
data "sotoon_iam_users" "all" {
  workspace_id = data.sotoon_workspace.mycompany.id
}

output "all_users" {
  description = "All users in the workspace"
  value       = data.sotoon_iam_users.all.users
}

# Retrieve all users of the "deployers" group
data "sotoon_iam_users" "deployers" {
  workspace_id = data.sotoon_workspace.mycompany.id
  group_id     = data.sotoon_iam_group.deployers.id
}

output "group_users" {
  description = "All users in the group"
  value       = data.sotoon_iam_users.deployers.users
}
