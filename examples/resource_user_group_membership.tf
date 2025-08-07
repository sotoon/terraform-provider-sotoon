locals {
  target_user_moein = one([
    for user in data.sotoon_iam_users.all.users : user if user.email == "moein.tavakoli@sotoon.ir"
  ])

  target_group_developer = one([
    for group in data.sotoon_iam_groups.all.groups : group if group.name == "developer"
  ])

  target_group_manager = one([
    for group in data.sotoon_iam_groups.all.groups : group if group.name == "manager"
  ])
}


resource "sotoon_iam_user_group_membership" "moein_dev_mng_membership" {
  user_id = local.target_user_moein.id
  group_ids = [local.target_group_developer.id , local.target_group_manager.id]

  depends_on = [ sotoon_iam_user.moein , sotoon_iam_group.developer ]
}
