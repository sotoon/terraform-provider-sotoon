data "sotoon_iam_service_users" "all" {}

output "service_users" {
  value = data.sotoon_iam_service_users.all.users
}
