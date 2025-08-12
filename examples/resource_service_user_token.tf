resource "sotoon_iam_service_user_token" "builder_token" {
  service_user_id = local.target_service_user_builder.id
}

output "service_user_token_binding" {
  value = sotoon_iam_service_user_token.builder_token.id
}
