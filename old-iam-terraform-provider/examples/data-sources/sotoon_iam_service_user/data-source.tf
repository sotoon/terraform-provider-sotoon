data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

data "sotoon_iam_service_user" "my_deployer" {
  workspace_id = data.sotoon_workspace.mycompany.id
  name         = "deployer"
}