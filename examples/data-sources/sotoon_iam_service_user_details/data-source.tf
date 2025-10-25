data "sotoon_iam_service_user_details" "builder" {
  service_user_id = "44444444-4444-4444-4444-444444444444"
}

output "service_user_name" {
  value = data.sotoon_iam_service_user_details.builder
}
