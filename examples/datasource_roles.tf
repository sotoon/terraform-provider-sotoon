data "sotoon_iam_roles" "all" {
  workspace_id = "4dc150c4-8e77-480f-98b6-821f8fc50227"
}

output "all_role_names" {
  description = "A list of all role names in the workspace."
  value       = data.sotoon_iam_roles.all.roles.*.name
}