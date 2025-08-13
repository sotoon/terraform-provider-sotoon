resource "sotoon_iam_user_role" "bind_moein_to_cdn" {
  role_id = "4b1d0837-8d72-42ad-bd04-5d19d357e12e"

  user_ids = [
    local.target_user_moein.id,
  ]
}

output "user_role_bind_id" {
  value = sotoon_iam_user_role.bind_moein_to_cdn.id
}
