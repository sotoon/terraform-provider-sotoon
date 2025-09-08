# Retrieve all IAM rules in the workspace
data "sotoon_iam_rules" "all" {
  workspace_id = var.workspace_id_target
}

# Output all rule names and their actions
output "all_rules" {
  description = "A list of all rules in the workspace"
  value = [
    for rule in data.sotoon_iam_rules.all.rules : {
      id      = rule.id
      name    = rule.name
      actions = rule.actions
      object  = rule.object
      deny    = rule.deny
    }
  ]
}

# Output all global rules
output "all_global_rules" {
  description = "A list of all global rules available to all workspaces"
  value = [
    for rule in data.sotoon_iam_rules.all.global_rules : {
      id      = rule.id
      name    = rule.name
      actions = rule.actions
      object  = rule.object
      deny    = rule.deny
    }
  ]
}
