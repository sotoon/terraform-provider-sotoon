# Contributing to the Sotoon Terraform Provider

Thank you for your interest in contributing! The Sotoon Provider is an open-source Terraform provider for managing [Sotoon](https://ocean.sotoon.ir/) resources. We welcome bug reports, feature requests, documentation improvements, and code contributions.

---

## Getting Started

### Prerequisites

* [Terraform](https://developer.hashicorp.com/terraform/install) >= **1.0**
* [Go](https://go.dev/doc/install) >= **1.18**
* Git & GitHub account

### Provider Configuration

Most development and testing requires a Sotoon account and API credentials.
Example provider block:

```
provider "sotoon" {
  api_token    = var.api_token
  workspace_id = var.workspace_id
  user_id      = var.user_id # Optional
}
```

Credentials can be passed via variables in [examples/terraform.tfvars](examples/terraform.tfvars).

---

## Development Workflow

1. Fork the repository.
2. Make your changes.
3. Add/update examples in the `examples/` directory if your change introduces or modifies resources.
4. Run formatting, linting, and tests ([see below](CONTRIBUTING.md#example-contributing-workflow)).
5. Open a Pull Request (PR) against `main`.

---

## Testing

Before submitting a PR, please run:

```bash
go fmt ./...
go vet ./...
go test ./...
```

### Acceptance Tests

To run against a real Sotoon environment:

```
  set credentials in examples/terraform.tfvars
```

Acceptance tests may incur real changes in your Sotoon account. Use a test workspace where possible.

---

## Issues & Pull Requests

* **Issues**:

  * Use clear titles and descriptions.
  * Include reproduction steps if reporting a bug.
  * Use labels: `bug`, `enhancement`, `docs`.

* **Pull Requests**:

  * Must reference an issue (e.g., “Fixes #123”).
  * Use conventional commit style in titles:

    * `feat: add support for X`
    * `fix: resolve crash when Y`
    * `docs: update example for Z`
    * `improve: rename variable i`
  * Include documentation and update `examples/` when relevant.

---

## Documentation

* All new resources and data sources must be documented under `docs/`.
* All new features must include an example under `examples/`.

---


## Security

Please **do not** report security vulnerabilities in public issues.
Instead, email to ```terraform@sotoon.ir```

---

## Example Contributing Workflow

This section illustrates how a new resource or data source is typically added to the Sotoon Terraform Provider.

1. **API Client Functions (`internal/client/client.go`)**

   * Define new functions to interact with Sotoon APIs.
   * These functions usually wrap calls from Sotoon’s Go SDKs (e.g., `sotoon/iamclient`).
   * Keep functions small and reusable.

2. **Resource or Data Source Implementation (`provider/*.go`)**

   * Create a new file for your resource or data source, following the naming conventions:

     * `resource_<name>.go` for resources
     * `data_source_<name>.go` for data sources
   * Use the client functions to implement CRUD (Create, Read, Update, Delete) operations.

3. **Provider Registration (`provider/provider.go`)**

   * Register your new resource or data source inside the `provider()` function, alongside the existing ones.

4. **Examples (`examples/`)**

   * Add or update an example Terraform configuration that demonstrates usage of the new or updated resource.
   * Ensure the example follows the file naming conventions used in this directory.

**Tip:** At each step, check and follow existing file naming conventions in the relevant directory for consistency.
---


## Acknowledgments

This project is inspired by the Terraform Provider ecosystem. Contributions are what make this provider possible — thank you for helping us improve Sotoon!

---