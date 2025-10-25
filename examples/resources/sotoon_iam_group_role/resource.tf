resource "sotoon_iam_group_role" "bind" {
  group_id = "33333333-3333-3333-3333-333333333333"
  role_ids = [
    "66666666-6666-6666-6666-666666666666",
    "66666666-6666-6666-6666-777777777777",
  ]
}

output "group_role_bind" {
  value = sotoon_iam_group_role.bind.id
}
