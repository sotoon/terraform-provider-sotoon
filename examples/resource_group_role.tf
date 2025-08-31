resource "sotoon_iam_group_role" "bind" {
  group_id = local.target_group_developer.id
  role_ids = [
    local.target_role_developer.id
  ]
}

output "group_role_bind" {
  value = sotoon_iam_group_role.bind.id
}
