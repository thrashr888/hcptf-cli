# Authentication Guide

## Quick Start

```bash
hcptf login
```

This prompts for an API token, validates it, and saves it to `~/.terraform.d/credentials.tfrc.json`. Both `hcptf` and `terraform` CLI share this credential file.

To generate a token, go to https://app.terraform.io/app/settings/tokens.

## Authentication Methods

The CLI checks for credentials in this order (first match wins):

| Priority | Method | Best for |
|----------|--------|----------|
| 1 | `TFE_TOKEN` env var | CI/CD pipelines |
| 2 | `HCPTF_TOKEN` env var | Separate hcptf credentials |
| 3 | `~/.hcptfrc` config file | Multiple instances, custom defaults |
| 4 | `~/.terraform.d/credentials.tfrc.json` | Shared with Terraform CLI |

### Environment Variables

```bash
# Standard Terraform variable - works with both terraform and hcptf
export TFE_TOKEN="your-token"

# hcptf-specific - use when you want a different token than terraform
export HCPTF_TOKEN="your-token"
```

### Configuration File

`~/.hcptfrc` (HCL format):

```hcl
credentials "app.terraform.io" {
  token = "your-token"
}

credentials "tfe.example.com" {
  token = "your-enterprise-token"
}

default_organization = "my-org"
output_format = "table"
```

### Terraform CLI Credentials

If you've already run `terraform login`, no additional setup is needed. `hcptf` reads `~/.terraform.d/credentials.tfrc.json` automatically.


Use `-show-token` to print the token for the target hostname without prompting:

```bash
hcptf login -show-token
# or for a custom hostname
hcptf login -hostname=tfe.example.com -show-token
```

Use `whoami` to verify token context and confirm the current user
(e.g., in scripts or after switching hosts):

```bash
hcptf whoami
hcptf whoami -output=json
```

## Multiple TFE Instances

Store credentials for multiple instances via `hcptf login`:

```bash
hcptf login                              # app.terraform.io (default)
hcptf login -hostname=tfe.example.com    # Terraform Enterprise
```

Switch between instances with environment variables:

```bash
# Default: app.terraform.io
hcptf workspace list -org=my-org

# Enterprise instance (preferred)
HCPTF_ADDRESS="https://tfe.example.com" hcptf workspace list -org=my-org

# Enterprise instance (legacy - also supported)
TFE_ADDRESS="https://tfe.example.com" hcptf workspace list -org=my-org
```

Note: `HCPTF_ADDRESS` takes precedence over `TFE_ADDRESS` for backward compatibility with Terraform Enterprise users.

## Logout

```bash
hcptf logout                             # Remove app.terraform.io credentials
hcptf logout -hostname=tfe.example.com   # Remove specific host credentials
```

This removes credentials from the local file only. To revoke tokens on the server, use the web UI.

## CI/CD Setup

Use the `TFE_TOKEN` environment variable. It works with both `terraform` and `hcptf`.

### GitHub Actions

```yaml
steps:
  - name: List workspaces
    env:
      TFE_TOKEN: ${{ secrets.TFE_TOKEN }}
    run: hcptf workspace list -org=my-org
```

### GitLab CI

```yaml
terraform:
  variables:
    TFE_TOKEN: ${TFE_TOKEN}
  script:
    - hcptf workspace list -org=my-org
```

## Token Types

| Type | Scope | Use case |
|------|-------|----------|
| User token | Your personal access | CLI usage, development |
| Team token | Team-level permissions | CI/CD, shared automation |
| Organization token | Organization-wide | Service accounts, broad automation |

## Troubleshooting

**"no authentication token found"** - No credentials configured. Run `hcptf login` or set `TFE_TOKEN`.

**"unauthorized" / 401** - Token is invalid or expired. Generate a new one at your TFE instance's token settings page.

**"permission denied" writing credentials** - Create the directory and retry:

```bash
mkdir -p ~/.terraform.d && chmod 755 ~/.terraform.d
hcptf login
```

**Works with terraform but not hcptf** - Check that `HCPTF_ADDRESS` points to the same host as your credentials. The token in the credentials file is keyed by hostname.
