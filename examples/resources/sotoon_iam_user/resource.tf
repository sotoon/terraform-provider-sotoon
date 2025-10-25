resource "sotoon_iam_user" "user" {
  email = "user@example.com"
}

output "user_id" {
  description = "The unique ID of the created user."
  value       = sotoon_iam_user.user.id
}
