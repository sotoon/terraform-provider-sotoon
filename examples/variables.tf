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

variable "user_id" {
  description = "user UUID"
  type        = string
  sensitive   = true
  nullable    = false
}

variable "ssh_public_key_path" {
  description = "Path to your SSH public key."
  type        = string
  default     = "~/.ssh/id_rsa.pub"
}

variable "workspace_id_target" {
  description = "The target workspace ID"
  type        = string
}

variable "user_email" {
  description = "Email address of the IAM user"
  type        = string
}