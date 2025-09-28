# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog, and this project adheres to Semantic Versioning.

## [0.1.0] - 2025-09-27

### Added
- First public release of the Sotoon Terraform Provider.
- Provider configuration with API token and workspace targeting (see `examples/terraform.tfvars.example`).
- Extensive examples under `examples/` and reference docs under `docs/`.

#### Supported Resources
- `sotoon_user`
- `sotoon_user_public_key`
- `sotoon_user_token`
- `sotoon_user_role`
- `sotoon_user_group_membership`
- `sotoon_group`
- `sotoon_group_role`
- `sotoon_role`
- `sotoon_service_user`
- `sotoon_service_user_public_key`
- `sotoon_service_user_role`
- `sotoon_service_user_token`
- `sotoon_service_user_group`

#### Supported Data Sources
- `sotoon_users`
- `sotoon_user`
- `sotoon_user_tokens`
- `sotoon_user_roles`
- `sotoon_user_public_key`
- `sotoon_groups`
- `sotoon_group_details`
- `sotoon_group_users`
- `sotoon_group_user_services`
- `sotoon_roles`
- `sotoon_rules`
- `sotoon_service_users`
- `sotoon_service_user_details`
- `sotoon_service_user_tokens`
- `sotoon_service_user_roles`
- `sotoon_service_user_public_keys`
- `sotoon_service_user_groups`
- `sotoon_service_user_groups_list`

### Compatibility
- Terraform >= 1.0
- Go >= 1.18 (for building the provider)

### Notes
- No breaking changes. This release focuses on core IAM resource CRUD and read capability, with comprehensive examples to get started quickly.
