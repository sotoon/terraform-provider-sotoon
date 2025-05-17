terraform {
  required_providers {
    sotoon = {
      source  = "terraform.sotoon.ir/sotoon/sotoon"
      version = "~> SOTOON_TERRAFORM_PROVIDER_REGISTRY"
    }
  }
}

provider "sotoon" {
  token = "xxxxxxxxxxxxxxxxxxxxxxxxxxx" # User Access Token
}
