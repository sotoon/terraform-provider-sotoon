terraform {
  required_providers {
    sotoon = {
      source  = "sotoon/sotoon"
      version = "0.1.6"
    }
  }
}

provider "sotoon" {
  api_token     = "your_sotoon_token_0000000000000000000000000000000000000000000000000000000000000000"
  workspace_id  = "11111111-1111-1111-1111-111111111111"
  user_id       = "22222222-2222-2222-2222-222222222222" # Optional
  should_log    = true # Optional
}


