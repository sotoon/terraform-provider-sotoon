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
