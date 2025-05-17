terraform {
  required_providers {
    sotoon = {
      source = "registry.terraform.io/cafebazaar/sotoon"
    }
  }
}

provider "sotoon" {
  host           = "https://api.sotoon.ir"
  zone           = "thr1"
  workspace_name = "cafebazaar"
  workspace_id   = "fee4bbf5-342d-4243-b4aa-f0bee508c39a"
}

# data "sotoon_subnet" "example" {
# }
#
# output "sotoon_subnet" {
#   value = data.sotoon_subnet.example
# }

# data "sotoon_externalip" "example" {
#  name = "auth-proxy"
# }

# output "sotoon_externalip" {
#  value = data.sotoon_externalip.example
# }

# data "sotoon_pvc" "example" {
# }
#
# output "sotoon_pvc" {
#  value = data.sotoon_pvc.example
# }


# resource "sotoon_subnet" "test" {
#   name = "test-tf"
#   cidr = "10.0.0.0/24"
# }
#
# output "subnet" {
#   value = sotoon_subnet.test
# }

# resource "sotoon_subnet" "test" {
#   name   = "test-tf"
#   cidr   = "10.0.0.0/24"
# #   routes = [
# #     {
# #       to          = "0.0.0.0/0"
# #       external_ip = "test"
# #     }
# #   ]
# }
#
# output "subnet" {
#   value = sotoon_subnet.test
# }
#
# resource "sotoon_pvc" "test" {
#   name = "test"
#   size = "11Gi"
#   tier = "standard"
# }
#
# output "pvc" {
#   value = sotoon_pvc.test
# }