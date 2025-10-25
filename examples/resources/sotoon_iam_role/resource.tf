resource "sotoon_iam_role" "basic" {
  name = "basic"
  description = "a simple role without no rules"
}

# Output the role IDs for reference
output "basic_role_id" {
  description = "The UUID of the basic role (no rules attached)"
  value       = sotoon_iam_role.basic.id
}

resource "sotoon_iam_role" "admin" {
  name = "admin"
  description = "a role with some rules"
  rules = [
    "77777777-7777-7777-7777-777777777777",
    "77777777-7777-7777-7777-888888888888",
  ]
}

# Output the role IDs for reference
output "admin_role_id" {
  description = "The UUID of the admin role"
  value       = sotoon_iam_role.admin.id
}
