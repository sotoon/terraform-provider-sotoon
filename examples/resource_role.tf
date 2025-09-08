# Create a new IAM role in the Sotoon workspace without any rules
resource "sotoon_iam_role" "developer" {
  name = "developer"
}

# Create a role with attached rules
resource "sotoon_iam_role" "admin" {
  name = "admin"
  
  # Attach specific rules to this role
  # These are rule UUIDs that must already exist in the system
  rules = [
    data.sotoon_iam_rules.all.rules[0].id,  # First rule from the workspace
    data.sotoon_iam_rules.all.global_rules[0].id  # First global rule
  ]
}

# Create a role with multiple rules
resource "sotoon_iam_role" "super_admin" {
  name = "super_admin"
  
  # Attach multiple rules to this role
  rules = [
    for rule in data.sotoon_iam_rules.all.rules : rule.id if contains(rule.actions, "create")
  ]
}

# Output the role IDs for reference
output "developer_role_id" {
  description = "The UUID of the developer role (no rules attached)"
  value       = sotoon_iam_role.developer.id
}

output "admin_role_id" {
  description = "The UUID of the admin role (with specific rules)"
  value       = sotoon_iam_role.admin.id
}

output "super_admin_role_id" {
  description = "The UUID of the super admin role (with multiple rules)"
  value       = sotoon_iam_role.super_admin.id
}