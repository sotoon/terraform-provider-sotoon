resource "sotoon_iam_service_user_public_key" "pk_by_file" {
  service_user_id = "44444444-4444-4444-4444-444444444444"
  title           = "key-by-file"
  public_key      = file("~/.ssh/id_rsa.pub")
}

output "service_user_public_key_binding" {
  value = sotoon_iam_service_user_public_key.pk_by_file.id
}

resource "sotoon_iam_service_user_public_key" "pk_by_key" {
  service_user_id = "44444444-4444-4444-4444-444444444444"
  title           = "key-by-key"
  public_key      = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDexamplefakekey000000000000000000000000000000000000000000000000000000000000 user@example.com"
}

output "service_user_public_key_binding" {
  value = sotoon_iam_service_user_public_key.pk_by_key.id
}
