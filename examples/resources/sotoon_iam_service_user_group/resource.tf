resource "sotoon_iam_service_user_group" "bind_builder_to_developer" {
  group_id = "33333333-3333-3333-3333-333333333333"
  service_user_ids = [
    "44444444-4444-4444-4444-444444444444",
  ]
}

output "group_user_bind" {
  value = sotoon_iam_service_user_group.bind_builder_to_developer.id
}


