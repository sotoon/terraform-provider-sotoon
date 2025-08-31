resource "sotoon_iam_service_user_group" "bind_builder_to_developer" {
  group_id = local.target_group_developer.id
  service_user_ids = [
    local.target_service_user_builder.id,
  ]
}

output "group_user_bind" {
  value = sotoon_iam_service_user_group.bind_builder_to_developer.id
}
