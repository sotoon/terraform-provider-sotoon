data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

data "sotoon_iam_group" "deployers" {
  workspace_id = data.sotoon_workspace.mycompany.id
  name         = "deployers"
}

data "sotoon_iam_service_user" "my_deployer" {
  workspace_id = data.sotoon_workspace.mycompany.id
  name         = "deployer"
}

# Joining `my_deployer` service-user to the deployers group
resource "sotoon_iam_service_user_group_binding" "my_service_user_is_deployer" {
  service_user_id = data.sotoon_iam_service_user.my_deployer.id
  workspace_id    = data.sotoon_workspace.mycompany.id
  group_id        = data.sotoon_iam_group.deployers.id
}
