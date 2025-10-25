resource "sotoon_iam_user_token" "me" {
  name      = "developer-token"
  expire_at = "2025-10-30T00:00:00Z"
}

output "new_user_token" {
  description = "The newly minted user token."
  value       = sotoon_iam_user_token.me.value
  sensitive   = true
}
