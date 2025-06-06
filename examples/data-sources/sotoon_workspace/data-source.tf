data "sotoon_workspace" "mycompany" {
  id = "ee6f89b5-e07c-42f1-9462-05cec9cd92d8" # Workspace ID
}

output "workspace_name" {
  description = "Name of my workspace"
  value       = data.sotoon_workspace.mycompany.name
}