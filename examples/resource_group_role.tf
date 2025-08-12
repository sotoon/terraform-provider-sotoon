resource "sotoon_iam_group_role" "bind" {
  group_id = "7f721307-3ee7-477f-8cb9-45010e230258"
  role_ids = [
    "4b1d0837-8d72-42ad-bd04-5d19d357e12e"
  ]
}

output "group_role_bind" {
  value = sotoon_iam_group_role.bind.id
}
