# Create a service user for CI/CD operations
resource "sotoon_iam_service_user" "builder" {
  name        = "builder"
  description = "Service user for CI/CD pipeline operations"
}

# Create another service user for monitoring
resource "sotoon_iam_service_user" "monitor" {
  name        = "monitor"
  description = "Service user for monitoring and alerting systems"
}

# Output the service user IDs for reference
output "builder_service_user_id" {
  description = "The UUID of the builder service user"
  value       = sotoon_iam_service_user.builder.id
}

output "monitor_service_user_id" {
  description = "The UUID of the monitor service user"
  value       = sotoon_iam_service_user.monitor.id
}
