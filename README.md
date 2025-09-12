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
    "sotoon/sotoon" = "path to go binary" # mac example: /Users/john/go/bin
  }
  direct {}
}
```

### Example Usage

This section shows a complete, runnable example.

1. Navigate to the examples directory:
```shell
cd examples/
```

2. Set environment variables:

```shell
cp terraform.tfvars.example terraform.tfvars
```
set values from [ocean panel](https://ocean.sotoon.ir/iam/users)

3. Initialize Terraform:

This command downloads any other required providers and sets up the backend.
```shell
terraform init
```

4. Apply the configuration:

This command will show you a plan of the resources to be created and prompt for confirmation.

```shell
terraform apply
```

5. Clean up resources:

When you are finished, destroy the created resources to avoid costs.
```shell
terraform destroy
```

## Contributing

We welcome contributions! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines and instructions.



# Running Tests

The provider includes **acceptance tests** that interact with real Sotoon APIs.  
These tests require valid credentials and environment variables.

## Prerequisites

- Go >= 1.18
- Valid Sotoon API credentials

## Setting Environment Variables

Before running the tests, set the following environment variables with your Sotoon account details:

```bash
export TF_ACC=1
export SOTOON_API_TOKEN="your-api-token"
export SOTOON_WORKSPACE_ID="your-workspace-id"
export SOTOON_USER_ID="your-user-id"
export SOTOON_API_HOST="https://bepa.sotoon.ir"
```

> 💡 Tip: You can put these into your shell profile (`~/.bashrc`, `~/.zshrc`) to avoid retyping.

## Running Acceptance Tests

To run all acceptance tests:

```bash
go test ./internal/... -v -tags=acceptance -run ^TestAcc -count=1
```

- `-v` → verbose output  
- `-tags=acceptance` → ensures only acceptance tests run  
- `-run ^TestAcc` → filters tests to only those starting with `TestAcc`  
- `-count=1` → disables test result caching  

## Example (all in one command)

If you prefer, you can run everything in one line without exporting variables separately:

```bash
TF_ACC=1 \
SOTOON_API_TOKEN="your-api-token" \
SOTOON_WORKSPACE_ID="your-workspace-id" \
SOTOON_USER_ID="your-user-id" \
SOTOON_API_HOST="https://bepa.sotoon.ir" \
go test ./internal/... -v -tags=acceptance -run ^TestAcc -count=1
```

⚠️ **Important:** These tests create and destroy real resources.  
Make sure you’re okay with changes happening in your Sotoon workspace.
