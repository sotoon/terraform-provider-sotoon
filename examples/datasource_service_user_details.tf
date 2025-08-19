data "sotoon_iam_service_user_details" "builder" {
  service_user_id = local.target_service_user_builder.id
}

output "service_user_name" {
  value = data.sotoon_iam_service_user_details.builder
}
