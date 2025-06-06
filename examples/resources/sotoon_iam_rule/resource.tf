data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

data "sotoon_service" "compute" {
  name = "compute"
}

resource "sotoon_iam_rule" "can_do_something" {
  name         = "can-do-something"
  workspace_id = data.sotoon_workspace.mycompany.id
  actions      = ["GET"]
  service      = data.sotoon_service.compute.id
  path         = "/path/to/some/resource/*"
  is_denial    = false
}
