terraform {
  required_providers {
    sotoon = {
      source  = "sotoon/sotoon"
      version = "1.0.0"
    }
  }
}


provider "sotoon" {
  api_token     = var.api_token
  workspace_id = var.workspace_id
}


