resource "sotoon_iam_group" "developer" {
  name = "developer"
  description = "this is user created by developers teams"
}

resource "sotoon_iam_group" "manager" {
  name = "manager"
}

output "group_id" {
  description = "The name of the created group of developer."
  value       = sotoon_iam_group.developer.name
}
