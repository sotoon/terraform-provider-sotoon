data "sotoon_iam_service_user_tokens" "this" {
  service_user_id = "44444444-4444-4444-4444-444444444444"
}

output "service_user_tokens" {
  value = data.sotoon_iam_service_user_tokens.this.tokens
}
