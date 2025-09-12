resource "sotoon_iam_user" "iam_user" {
  email = var.user_email
}

output "user_id" {
  description = "The unique ID of the created user."
  value       = sotoon_iam_user.iam_user.id
}
