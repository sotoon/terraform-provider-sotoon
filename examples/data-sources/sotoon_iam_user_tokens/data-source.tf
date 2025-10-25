data "sotoon_iam_user_tokens" "me" {}

output "existing_user_token_ids" {
  description = "All token UUIDs for your account."
  value       = data.sotoon_iam_user_tokens.me.tokens[*].id
}
