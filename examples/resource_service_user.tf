resource "sotoon_iam_service_user" "builder" {
  name = "builder"
  description = "this is builder in cicd"
}

output "service_user_id" {
  value = sotoon_iam_service_user.builder.id
}
