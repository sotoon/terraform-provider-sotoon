data "sotoon_iam_user_public_keys" "me" {}

output "existing_public_key_ids" {
  description = "All SSH public-key UUIDs for your account."
  value       = data.sotoon_iam_user_public_keys.me.public_keys[*].id
}
