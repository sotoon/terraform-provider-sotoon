resource "sotoon_iam_user_role" "bind_user_to_role" {
  role_id = "77777777-7777-7777-7777-777777777777"

  user_ids = [
    "44444444-4444-4444-4444-444444444444",
  ]
}

output "user_role_bind_id" {
  value = sotoon_iam_user_role.bind_user_to_role.id
}
