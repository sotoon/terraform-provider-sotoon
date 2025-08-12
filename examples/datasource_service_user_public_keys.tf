data "sotoon_iam_service_user_public_keys" "this" {
  service_user_id = "29a54954-3f38-4a5e-9e51-1ea582ab90f2"
}

output "service_user_public_keys" {
  value = data.sotoon_iam_service_user_public_keys.this.public_keys
}
