data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

data "sotoon_iam_user" "me" {
  workspace_id = data.sotoon_workspace.mycompany.id
  email        = "doe@mycompany.com"
}
