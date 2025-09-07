# Create a new IAM role in the Sotoon workspace
resource "sotoon_iam_role" "developer" {
  name = "developer"
}

# Create another role with a different name
resource "sotoon_iam_role" "admin" {
  name = "admin"
}

# Output the role IDs for reference
output "developer_role_id" {
  description = "The UUID of the developer role"
  value       = sotoon_iam_role.developer.id
}

output "admin_role_id" {
  description = "The UUID of the admin role"
  value       = sotoon_iam_role.admin.id
}