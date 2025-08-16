resource "sotoon_iam_user_role" "bind_moein_to_developer" {
  role_id = local.target_role_developer.id

  user_ids = [
    local.target_user_moein.id,
  ]
}

output "user_role_bind_id" {
  value = sotoon_iam_user_role.bind_moein_to_developer.id
}
