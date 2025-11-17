resource "sotoon_iam_service_user_role" "bind_service_user_to_role" {
  role_id = "77777777-7777-7777-7777-777777777777"
  service_user_ids = [
    "44444444-4444-4444-4444-444444444444",
  ]
  items = {
    "key1" = "value1",
    "key2" = "value2",
  }
}

output "service_user_role_binding" {
  value = sotoon_iam_service_user_role.bind_service_user_to_role.id
}
