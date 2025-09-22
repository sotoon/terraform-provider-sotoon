# Terraform Provider Sotoon
The Terraform Sotoon provider is a plugin for Terraform that allows for the full lifecycle management of Sotoon resources.

## Requirements

- Terraform >= 1.0
- Go >= 1.18 (to build the provider plugin)

## Getting Started

1. Clone the repository:

```shell
git clone https://github.com/sotoon/terraform-provider-sotoon
cd terraform-provider-sotoon
```

2. Build the provider:
   
Run the following command to build the provider binary.
This will create the terraform-provider-sotoon executable in the project's root directory.

```shell
go build -o terraform-provider-sotoon
```

#### Installing the Provider for Development

For development purposes, you can instruct Terraform to use your locally-built provider instead of downloading one from a registry.

1. Install with Make:

The provided Makefile simplifies the process. It builds the binary and copies it to the correct plugins directory for your user.

```shell
make install
```

2. Manual Installation (Alternative):
You can configure Terraform manually. Create or edit the Terraform CLI configuration file (~/.terraformrc or %APPDATA%\terraform.rc on Windows).

```shell
provider_installation {
  dev_overrides {
    "sotoon/sotoon" = "<path-to-your-built-binary>"
  }
  direct {}
}
```
Replace `<path-to-your-built-binary>` with the directory containing your built `terraform-provider-sotoon` binary.

### Example Usage

This section shows a complete, runnable example.

1. Navigate to the examples directory:
```shell
cd examples/
```

2. Initialize Terraform:

This command downloads any other required providers and sets up the backend.
```shell
terraform init
```

3. Apply the configuration:

This command will show you a plan of the resources to be created and prompt for confirmation.

```shell
terraform apply
```

4. Clean up resources:

When you are finished, destroy the created resources to avoid costs.
```shell
terraform destroy
```

## Running Tests

To run unit tests locally:

```shell
go test ./internal/provider/... -tags=unit -v -run '^TestUnit'

## Example `terraform.tfvars`

Copy the example file before customizing:

```shell
cp examples/terraform.tfvars.example examples/terraform.tfvars
```

Below is an example of how to configure your `terraform.tfvars` file:

```hcl
api_token     = "xxxxxxxxxxxxxxxxxxxxxxxxxxx"
workspace_id  = "your-workspace-id"
user_id       = "your user id"

# Special variable used in users, groups, etc.
workspace_id_target = "your-target-workspace-id"
```

## Contributing

We welcome contributions! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines and instructions.