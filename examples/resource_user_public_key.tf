resource "sotoon_iam_user_public_key" "myself" {
  title      = "pc-ssh"
  public_key = file(var.ssh_public_key_path)
}

resource "sotoon_iam_user_public_key" "sotoon"{
  title = "sotoon_pub"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAA ..."
  # Default Key type is RSA
}

resource "sotoon_iam_user_public_key" "sotoon_ed25519"{
  title = "sotoon-pc"
  public_key = file("~/.ssh/id_ed25519.pub")
  key_type = "ssh-ed25519"
}

output "new_public_key_id" {
  description = "The UUID of the created SSH public key."
  value       = sotoon_iam_user_public_key.myself.id
}



