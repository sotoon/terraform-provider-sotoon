resource "sotoon_iam_service_user_role" "bind_builder_to_developer" {
  role_id = local.target_role_developer.id
  service_user_ids = [
    local.target_service_user_builder.id,
  ]
}

output "service_user_role_binding" {
  value = sotoon_iam_service_user_role.bind_builder_to_developer.id
}
