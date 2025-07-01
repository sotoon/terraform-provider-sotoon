variable "api_token" {
  description = "API token"
  type        = string
  sensitive   = true
  nullable    = false
}

variable "workspace_id" {
  description = "Workspace UUID"
  type        = string
  sensitive   = true
  nullable    = false
}

