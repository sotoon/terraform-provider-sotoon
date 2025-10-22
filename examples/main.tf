terraform {
  required_providers {
    sotoon = {
      source  = "sotoon/sotoon"
      version = "0.1.3"
    }
  }
}

provider "sotoon" {
  api_token     = var.api_token
  workspace_id  = var.workspace_id
  user_id       = var.user_id # Optional
  should_log    = true # Optional
}


