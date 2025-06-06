data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

data "sotoon_iam_group" "deployers" {
  workspace_id = data.sotoon_workspace.mycompany.id
  name         = "deployers"
}

data "sotoon_iam_role" "compute_viewer" {
  name = "compute-viewer"
}

resource "sotoon_iam_role_group_binding" "deployers_are_compute_viewers" {
  group_id     = data.sotoon_iam_group.deployers.id
  workspace_id = data.sotoon_workspace.mycompany.id
  role_id      = data.sotoon_iam_role.compute_viewer.id
  items = {
    "zone" : "neda",
    "namespace" : "mycompany"
  }
}
