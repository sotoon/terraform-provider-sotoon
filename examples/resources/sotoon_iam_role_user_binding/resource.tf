data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

data "sotoon_iam_user" "john" {
  workspace_id = data.sotoon_workspace.mycompany.id
  email        = "john.doe@sotoon.ir"
}

data "sotoon_iam_role" "compute_viewer" {
  name = "compute-viewer"
}

resource "sotoon_iam_role_user_binding" "john_is_compute_viewer" {
  user_id      = data.sotoon_iam_user.john.id
  workspace_id = data.sotoon_workspace.mycompany.id
  role_id      = data.sotoon_iam_role.compute_viewer.id
  items = {
    "zone" : "neda",
    "namespace" : "mycompany"
  }
}
