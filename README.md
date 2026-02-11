# HCP Terraform CLI (`hcptf`)

A Go CLI for managing HCP Terraform resources. Built with `mitchellh/cli` and `hashicorp/go-tfe`.

## Installation

```bash
go build -o hcptf .

# Optional: add to PATH
sudo mv hcptf /usr/local/bin/
```

## Authentication

The fastest way to get started:

```bash
hcptf login
```

This prompts for an API token, validates it, and stores it in `~/.terraform.d/credentials.tfrc.json` (shared with Terraform CLI).

If you've already run `terraform login`, no setup is needed - `hcptf` reads those credentials automatically.

Other authentication methods (checked in order):

1. `TFE_TOKEN` environment variable
2. `HCPTF_TOKEN` environment variable
3. `~/.hcptfrc` configuration file (HCL format)
4. `~/.terraform.d/credentials.tfrc.json`

See [docs/AUTH_GUIDE.md](docs/AUTH_GUIDE.md) for details on CI/CD setup, multiple TFE instances, and troubleshooting.

## Configuration

Optional config file at `~/.hcptfrc`:

```hcl
credentials "app.terraform.io" {
  token = "your-token-here"
}

default_organization = "my-org"
output_format = "table"  # or "json"
```

Override the API endpoint with `HCPTF_ADDRESS`:

```bash
export HCPTF_ADDRESS="https://tfe.example.com"
```

## Usage

### Traditional Command Style

```bash
# List workspaces
hcptf workspace list -org=my-org

# Create workspace and trigger a run
hcptf workspace create -org=my-org -name=staging -auto-apply=false
hcptf run create -org=my-org -workspace=staging -message="Deploy changes"

# Check run status and apply
hcptf run show -id=run-abc123
hcptf run apply -id=run-abc123 -comment="Approved"

# Manage variables
hcptf variable create -org=my-org -workspace=staging -key=region -value=us-east-1
hcptf variable create -org=my-org -workspace=staging \
  -key=AWS_SECRET_KEY -value=secret -category=env -sensitive

# JSON output for scripting
hcptf workspace list -org=my-org -output=json
```

### URL-Style Navigation

For convenience, you can use a URL-like path syntax to access resources:

```bash
# Show organization details
hcptf my-org

# List workspaces in an org
hcptf my-org workspaces

# Show workspace details
hcptf my-org my-workspace

# List runs for a workspace
hcptf my-org my-workspace runs
hcptf my-org my-workspace runs list

# Show a specific run
hcptf my-org my-workspace runs run-abc123

# Apply a run
hcptf my-org my-workspace runs run-abc123 apply

# List workspace variables
hcptf my-org my-workspace variables

# View state
hcptf my-org my-workspace state
hcptf my-org my-workspace state outputs

# Other org resources
hcptf my-org projects
hcptf my-org teams
hcptf my-org policies
```

Both styles work interchangeably - use whichever you prefer!

## Commands

229 commands across 59 resource types. See [docs/COMMANDS.md](docs/COMMANDS.md) for the full reference.

| Group | Commands | Description |
|-------|----------|-------------|
| `login` / `logout` | 2 | Credential management |
| `account` | 3 | User account CRUD |
| `workspace` | 5 | Workspace management |
| `run` | 6 | Run lifecycle |
| `organization` | 5 | Organization management |
| `variable` | 4 | Workspace variables |
| `team` | 6 | Teams and membership |
| `project` | 5 | Project organization |
| `state` | 3 | State versions and outputs |
| `policy` | 5 | Sentinel policies |
| `policyset` | 7 | Policy set management |
| `sshkey` | 5 | SSH keys for VCS |
| `notification` | 6 | Run notifications |
| `variableset` | 10 | Reusable variable sets |
| `agentpool` | 8 | Self-hosted agent pools |
| `runtask` | 7 | Run task integrations |
| `oauthclient` | 5 | VCS OAuth clients |
| `oauthtoken` | 3 | OAuth tokens |
| `runtrigger` | 4 | Workspace orchestration |
| `plan` | 2 | Plan details and logs |
| `apply` | 2 | Apply details and logs |
| `configversion` | 4 | Configuration versions |
| `teamaccess` | 5 | Team workspace permissions |
| `projectteamaccess` | 5 | Team project permissions |
| `registrymodule` | 6 | Private registry modules |
| `registryprovider` | 4 | Private registry providers |
| `registryproviderversion` | 3 | Provider versions |
| `registryproviderplatform` | 3 | Provider platforms |
| `gpgkey` | 5 | GPG keys for providers |
| `stack` | 5 | Terraform Stacks |
| `stackconfiguration` | 5 | Stack configurations |
| `stackdeployment` | 3 | Stack deployments |
| `stackstate` | 2 | Stack states |
| `audittrail` | 2 | Audit trail events |
| `audittrailtoken` | 4 | Audit trail tokens |
| `version` | 1 | CLI version |

### Common flags

| Flag | Alias | Description |
|------|-------|-------------|
| `-organization` | `-org` | Organization name |
| `-output` | | `table` (default) or `json` |
| `-force` | | Skip confirmation prompts |

### Help

```bash
hcptf --help
hcptf workspace --help
hcptf workspace create --help
```

## Project Structure

```
hcptf-cli/
├── main.go                  # Entry point
├── command/                 # CLI commands
│   ├── commands.go          # Command registry
│   ├── meta.go              # Shared command base
│   └── <resource>_<action>.go
├── internal/
│   ├── client/              # go-tfe API client wrapper
│   ├── config/              # Configuration loading
│   └── output/              # Table/JSON formatting
├── docs/                    # Documentation
└── go.mod
```

## Development

```bash
go build -o hcptf .
go test ./...
./hcptf --help
```

### Dependencies

- `github.com/hashicorp/go-tfe` - HCP Terraform API client
- `github.com/mitchellh/cli` - CLI framework
- `github.com/hashicorp/hcl/v2` - Configuration parsing
- `github.com/olekukonko/tablewriter` - Table output
- `github.com/fatih/color` - Terminal colors

### Exit codes

- `0` - Success
- `1` - API or runtime error
- `2` - Usage or validation error

## License

Designed for use with HCP Terraform and Terraform Enterprise.
