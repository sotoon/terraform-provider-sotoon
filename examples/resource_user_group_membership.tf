resource "sotoon_iam_user_group_membership" "dev_group_members" {
  group_id = local.target_group_developer.id
  user_ids = [
    local.target_user_moein.id,
  ]

  depends_on = [
    sotoon_iam_group.developer,
    sotoon_iam_user.moein,
  ]
}
