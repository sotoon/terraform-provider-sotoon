resource "sotoon_iam_service_user_token" "builder_token" {
  service_user_id =  "44444444-4444-4444-4444-444444444444"
}

output "service_user_token_value" {
  description = "Service user token"
  value       = sotoon_iam_service_user_token.builder_token.value
  sensitive   = true
}
