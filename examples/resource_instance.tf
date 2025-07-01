resource "sotoon_instance" "my_vm" {
  name  = "test-instance-terraform-1"
  type  = "c1d30-1"
  image = "ubuntu-22.04"
}


resource "sotoon_instance" "zurvun_vm" {
  name  = "zurvun-terraform-0"
  type  = "g1d-1"
  image = "oracle-linux-9.5"
  iam_enabled = true 

}

