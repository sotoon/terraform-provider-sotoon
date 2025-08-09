resource "sotoon_iam_user_group_membership" "moein_dev_mng_membership" {
  user_id = local.target_user_moein.id
  group_ids = [local.target_group_developer.id , local.target_group_manager.id]

  depends_on = [ sotoon_iam_user.moein , sotoon_iam_group.developer ]
}
