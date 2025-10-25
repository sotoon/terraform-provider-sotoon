resource "sotoon_iam_user_group_membership" "dev_group_members" {
  group_id = "33333333-3333-3333-3333-333333333333"
  user_ids = [
    "44444444-4444-4444-4444-444444444444",
  ]
}
