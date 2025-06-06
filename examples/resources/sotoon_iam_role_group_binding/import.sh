# Import ID for this resource is "<role-id>:<group-id>:<workspace-id>".
# Replace <role-id>, <group-id> and <workspace-id> with real values
terraform import sotoon_iam_role_group_binding.deployers_are_compute_viewers "<role-id>:<group-id>:<workspace-id>"

# For example:
terraform import sotoon_iam_role_group_binding.deployers_are_compute_viewers \
    "b8c133a4-a060-4906-8654-57988dbdf098:34f57a2f-6e4d-4ded-9025-ff00911d3313:ee6f89b5-e07c-42f1-9462-05cec9cd92d8"
