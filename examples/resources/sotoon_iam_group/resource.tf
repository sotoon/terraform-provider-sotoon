data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

# Create `deployers` group in the `mycompany` workspace
resource "sotoon_iam_group" "deployers" {
  name         = "deployers"
  workspace_id = data.sotoon_workspace.mycompany.id
}
