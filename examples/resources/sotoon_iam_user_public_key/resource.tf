resource "sotoon_iam_user_public_key" "by_file" {
  title      = "pc-ssh"
  public_key = file("~/.ssh/id_rsa.pub")
}

resource "sotoon_iam_user_public_key" "by_key"{
  title = "sotoon_pub"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAA ..."
  # Default Key type is RSA
}

output "new_public_key_id_by_file" {
  description = "The UUID of the created SSH public key."
  value       = sotoon_iam_user_public_key.by_file.id
}

output "new_public_key_id_by_key" {
  description = "The UUID of the created SSH public key."
  value       = sotoon_iam_user_public_key.by_key.id
}

