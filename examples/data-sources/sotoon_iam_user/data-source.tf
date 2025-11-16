
data "sotoon_iam_user" "by_email" {
  workspace_id = "11111111-1111-1111-1111-111111111111"
  email        = "user@example.com"
}

output "user_by_email" {
  value = data.sotoon_iam_user.by_email.uuid
}

data "sotoon_iam_user" "by_uuid" {
  workspace_id = "11111111-1111-1111-1111-111111111111"
  uuid         = "22222222-2222-2222-2222-222222222222"
}

output "user_by_uuid" {
  value = data.sotoon_iam_user.by_uuid.email
}
