data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

# Get globally predefined workspace-admin role 
data "sotoon_iam_role" "workspace_admin" {
  name = "workspace-admin"
}

# Get user defined custom role
data "sotoon_iam_role" "reader" {
  name         = "my-user-can-read"
  workspace_id = data.sotoon_workspace.mycompany.id
}
