resource "sotoon_iam_user" "moein" {
  email = "moein.tavakoli-test4@zurvun.com"
}

output "user_id" {
  description = "The unique ID of the created user."
  value       = sotoon_iam_user.moein.id
}
