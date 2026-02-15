# HCP Terraform CLI (`hcptf`)

A Go CLI for managing HCP Terraform resources. Built with `mitchellh/cli` and `hashicorp/go-tfe`.

## Installation

### Download Pre-built Binary (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/thrashr888/hcptf-cli/releases):

```bash
# Example: macOS (Apple Silicon)
curl -LO https://github.com/thrashr888/hcptf-cli/releases/latest/download/hcptf_VERSION_darwin_arm64.tar.gz
tar -xzf hcptf_VERSION_darwin_arm64.tar.gz
sudo mv hcptf /usr/local/bin/

# Verify installation
hcptf version
```

Available platforms: Linux (amd64, arm64), macOS (Intel, Apple Silicon), Windows, FreeBSD.

See [docs/RELEASING.md](docs/RELEASING.md) for checksum verification.

### Build from Source

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

Override the API endpoint:

```bash
export HCPTF_ADDRESS="https://tfe.example.com"  # Preferred
export TFE_ADDRESS="https://tfe.example.com"    # Legacy support
```

Note: `HCPTF_ADDRESS` takes precedence over `TFE_ADDRESS` for compatibility.

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

# Registry commands (hierarchical namespace)
hcptf registry module list -org=my-org
hcptf registry provider create -org=my-org -name=custom-provider
hcptf registry provider version create -org=my-org -name=aws -version=3.1.1 -key-id=GPG_KEY_ID

# Stack commands (hierarchical namespace)
hcptf stack list -org=my-org -project=prj-abc123
hcptf stack configuration list -stack-id=stk-abc123
hcptf stack deployment create -stack-id=stk-abc123

# Explorer API (query resources across organization)
hcptf explorer query -org=my-org -type=providers -sort=-version
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

### Hierarchical Command Namespaces

Registry and stack commands use a hierarchical namespace structure for better organization:

```bash
# Registry commands
hcptf registry                                    # Show all registry commands
hcptf registry module list -org=my-org           # List modules
hcptf registry provider create -org=my-org       # Create provider
hcptf registry provider version create ...        # Create provider version
hcptf registry provider platform create ...       # Add platform binary

# Stack commands
hcptf stack                                       # Show all stack commands
hcptf stack list -org=my-org                     # List stacks
hcptf stack configuration list -stack-id=stk-123  # List configurations
hcptf stack deployment create -stack-id=stk-123   # Create deployment
hcptf stack state list -stack-id=stk-123         # List state versions
```

## Commands

100+ commands across 50+ resource types organized in hierarchical namespaces.

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
| `registry` | 1 | Private registry parent |
| `registry module` | 6 | Private registry modules |
| `registry provider` | 4 | Private registry providers |
| `registry provider version` | 3 | Provider versions |
| `registry provider platform` | 3 | Provider platforms |
| `gpgkey` | 5 | GPG keys for providers |
| `stack` | 6 | Terraform Stacks management |
| `stack configuration` | 5 | Stack configurations |
| `stack deployment` | 3 | Stack deployments |
| `stack state` | 2 | Stack states |
| `audittrail` | 2 | Audit trail events |
| `audittrailtoken` | 4 | Audit trail tokens |
| `explorer` | 1 | Query resources across org |
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

## Agent Skills

This project includes [Agent Skills](https://agentskills.io/) that help AI agents use the CLI effectively. Skills are automatically discovered by compatible agents (Claude Code, Cursor, GitHub Copilot, etc.).

**Available skills:**
- **hcptf-cli**: Comprehensive guide covering authentication, commands, workflows, and best practices

The `.skills/` directory contains structured instructions that agents can load to:
- Understand hierarchical command structure
- Learn common workflows (workspace creation, deployments, etc.)
- Follow best practices for automation
- Handle errors and troubleshooting

See [.skills/README.md](.skills/README.md) for details.

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
