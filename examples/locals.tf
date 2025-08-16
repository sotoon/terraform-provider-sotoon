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
  target_service_user_builder = one([
    for su in data.sotoon_iam_service_users.all.users : su if su.name == "builder"
  ])
  target_role_developer = one([
    for role in data.sotoon_iam_roles.all.roles : role if role.name == "developer"
  ])
}
