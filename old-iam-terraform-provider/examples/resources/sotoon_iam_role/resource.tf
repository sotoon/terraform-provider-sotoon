data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

data "sotoon_service" "compute" {
  name = "compute"
}

# Create `can-do-something` rule in `mycompany` workspace
resource "sotoon_iam_rule" "can_do_something" {
  name         = "can-do-something"
  workspace_id = data.sotoon_workspace.mycompany.id
  actions      = ["GET"]
  service      = data.sotoon_service.compute.id
  path         = "/path/to/some/resource/*"
  is_denial    = false
}

# Load `can-do-another-thing` global rule
data "sotoon_iam_rule" "can_do_another_thing" {
  name = "can-do-another-thing"
}

# Create `my-role` role that two `can-do-another-thing` and `can-do-something` rules included.
resource "sotoon_iam_role" "my_role" {
  name         = "my-role"
  workspace_id = data.sotoon_workspace.mycompany.id
  rules = [
    { id = sotoon_iam_rule.can_do_something.id },
    { id = data.sotoon_iam_rule.can_do_another_thing.id },
  ]
}
