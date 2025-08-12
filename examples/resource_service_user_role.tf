resource "sotoon_iam_service_user_role" "bind_builder_to_compute" {
  role_id = "72ddc02f-57e8-42c7-8bb3-cc1736474272"

  service_user_ids = [
    local.target_service_user_builder.id,
  ]
}

output "service_user_role_binding" {
  value = sotoon_iam_service_user_role.bind_builder_to_compute.id
}
