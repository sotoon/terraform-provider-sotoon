data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

data "sotoon_iam_group" "deployers" {
  workspace_id = data.sotoon_workspace.mycompany.id
  name         = "deployers"
}

data "sotoon_iam_user" "john" {
  workspace_id = data.sotoon_workspace.mycompany.id
  email        = "john.doe@sotoon.ir"
}

# Joining John to the deployers group
resource "sotoon_iam_user_group_binding" "john_is_deployer" {
  user_id      = data.sotoon_iam_user.john.id
  workspace_id = data.sotoon_workspace.mycompany.id
  group_id     = data.sotoon_iam_group.deployers.id
}
