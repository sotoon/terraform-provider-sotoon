resource "sotoon_iam_service_user_public_key" "builder_pk" {
  service_user_id = local.target_service_user_builder.id
  title           = "builder-key"
  public_key      = file("~/.ssh/id_rsa.pub")
}

output "service_user_public_key_binding" {
  value = sotoon_iam_service_user_public_key.builder_pk.id
}
