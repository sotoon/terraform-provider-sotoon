data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

data "sotoon_iam_service_user" "my_deployer" {
  workspace_id = data.sotoon_workspace.mycompany.id
  name         = "deployer"
}

data "sotoon_iam_role" "compute_viewer" {
  name = "compute-viewer"
}

resource "sotoon_iam_role_service_user_binding" "mydeployer_is_compute_viewer" {
  service_user_id = data.sotoon_iam_service_user.my_deployer.id
  workspace_id    = data.sotoon_workspace.mycompany.id
  role_id         = data.sotoon_iam_role.compute_viewer.id
  items = {
    "zone" : "neda",
    "namespace" : "mycompany"
  }
}
