resource "sotoon_iam_user" "moein" {
  name     = "Moein Tavakoli"
  email    = "moein.tavakoli@sotoon.ir"
  password = "S0me_R@nd0m_Pa$$w0rd"
}

output "user_id" {
  description = "The unique ID of the created user."
  value       = sotoon_iam_user.moein.id
}
