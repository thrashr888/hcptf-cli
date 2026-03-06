---
name: hcptf-cli
description: Manage HCP Terraform resources using the hcptf command-line tool. Use when working with HCP Terraform, Terraform Cloud, Terraform Enterprise, workspaces, runs, organizations, stacks, registry modules/providers, or infrastructure automation tasks.
---

# HCP Terraform CLI (hcptf)

Use the CLI with workspace/run hierarchical commands and URL-style navigation.

## When to Use This Skill

Use this skill when:

- Managing HCP Terraform workspaces, runs, or organizations
- Working with Terraform Stacks or private registry resources
- Automating infrastructure deployments
- Querying resource state across organizations
- User mentions "Terraform Cloud", "HCP Terraform", or "TFE"

## Authentication

The CLI checks for credentials in this order:

1. `TFE_TOKEN` environment variable
2. `HCPTF_TOKEN` environment variable
3. `~/.hcptfrc` configuration file
4. `~/.terraform.d/credentials.tfrc.json` (shared with Terraform CLI)

### Quick Setup

```bash
# Interactive login (recommended)
hcptf login

# Or set environment variable
export TFE_TOKEN="your-token-here"
```

## Command Structure

The CLI uses a **hierarchical namespace** structure for organization:

### Hierarchical Commands

```bash
# Registry commands
hcptf registry                                    # Show all registry commands
hcptf registry module list -org=my-org
hcptf registry provider create -org=my-org -name=custom
hcptf registry provider version create -org=my-org -name=aws -version=3.1.1 -key-id=GPG_KEY

# Stack commands
hcptf stack                                       # Show all stack commands
hcptf stack list -org=my-org -project=prj-123
hcptf stack configuration list -stack-id=stk-123
hcptf stack deployment create -stack-id=stk-123
hcptf stack state list -stack-id=stk-123

# Explorer API (query across org)
hcptf explorer query -org=my-org -type=providers -sort=-version
```

### Traditional Commands

```bash
# Workspace operations
hcptf workspace list -org=my-org
hcptf workspace create -org=my-org -name=staging
hcptf workspace read -org=my-org -name=staging

# Run operations
hcptf workspace run list -org=my-org -name=staging
hcptf workspace run create -org=my-org -name=staging -message="Deploy"
hcptf workspace run show -id=run-abc123
hcptf workspace run apply -id=run-abc123
hcptf workspace run cancel -id=run-abc123
hcptf workspace run discard -id=run-abc123
hcptf workspace run assessmentresult list -org=my-org -name=staging
hcptf workspace run assessmentresult read -id=asmtres-abc123 -show-drift

# Variables
hcptf variable create -org=my-org -workspace=staging -key=region -value=us-east-1
hcptf variable create -org=my-org -workspace=staging \
  -key=AWS_SECRET -value=secret -category=env -sensitive
```

### URL-Style Navigation

For interactive exploration, use path-like syntax:

```bash
# Organization level
hcptf my-org                    # Show org details
hcptf my-org workspaces         # List workspaces
hcptf my-org projects           # List projects
hcptf my-org teams              # List teams

# Workspace level
hcptf my-org my-workspace       # Show workspace details
hcptf my-org my-workspace runs  # List runs
hcptf my-org my-workspace assessments  # List assessment results
hcptf my-org my-workspace variables
hcptf my-org my-workspace state

# Run operations
hcptf my-org my-workspace runs run-abc123        # Show run
hcptf my-org my-workspace runs run-abc123 apply  # Apply run
```

## Common Workflows

### Create Workspace and Deploy

```bash
# Create workspace
hcptf workspace create -org=my-org -name=production \
  -auto-apply=false \
  -terraform-version=1.6.0

# Set variables
hcptf variable create -org=my-org -workspace=production \
  -key=environment -value=prod

# Trigger run
hcptf workspace run create -org=my-org -name=production \
  -message="Initial deployment"

# Check status
hcptf workspace run show -id=run-abc123

# Apply when ready
hcptf workspace run apply -id=run-abc123 -comment="Approved by team"
```

### Manage Private Registry

```bash
# List modules
hcptf registry module list -org=my-org

# Create provider
hcptf registry provider create -org=my-org \
  -name=custom-aws \
  -namespace=my-org

# Add provider version
hcptf registry provider version create -org=my-org \
  -name=custom-aws \
  -version=1.0.0 \
  -key-id=ABC123

# Add platform binary
hcptf registry provider platform create -org=my-org \
  -name=custom-aws \
  -version=1.0.0 \
  -os=linux \
  -arch=amd64 \
  -filename=terraform-provider-custom-aws_1.0.0_linux_amd64.zip
```

### Work with Stacks

```bash
# List stacks
hcptf stack list -org=my-org -project=prj-abc123

# Create stack
hcptf stack create -org=my-org \
  -project=prj-abc123 \
  -name=production-stack

# List configurations
hcptf stack configuration list -stack-id=stk-abc123

# Trigger deployment
hcptf stack deployment create -stack-id=stk-abc123

# Check deployment status
hcptf stack deployment read -deployment-id=dep-abc123

# View state
hcptf stack state list -stack-id=stk-abc123
```

### Query Resources with Explorer

```bash
# Find providers across organization
hcptf explorer query -org=my-org \
  -type=providers \
  -sort=-version

# Query workspaces
hcptf explorer query -org=my-org \
  -type=workspaces \
  -filter="name:production"

# Export to CSV
hcptf explorer query -org=my-org \
  -type=tf_versions \
  -export-csv \
  -output=versions.csv
```

## Agent-Friendly Features

The CLI is designed to work well when driven by AI agents. **Always use these features** when automating with hcptf.

### Schema Introspection

Before constructing commands, discover available flags as structured JSON:

```bash
# Get machine-readable flag schema for any command
hcptf schema workspace create
hcptf schema variable create
hcptf schema run apply
```

Returns JSON with flag names, types, required status, aliases, and defaults. **Always use `hcptf schema` instead of parsing `--help` text.**

### Dry Run (validate without executing)

**Always use `-dry-run` before mutations** to validate inputs without making API calls:

```bash
# Preview what would happen
hcptf workspace delete -org=my-org -name=staging -dry-run
hcptf variable create -org=my-org -workspace=prod -key=region -value=us-east-1 -dry-run
hcptf run apply -id=run-abc123 -dry-run
```

Returns JSON describing the planned action. No API call is made.

### JSON Input for Mutations

Pass full API payloads instead of individual flags using `-json-input`:

```bash
# Inline JSON
hcptf variable create -org=my-org -workspace=prod \
  -json-input='{"key":"region","value":"us-east-1","category":"terraform"}'

# From file
hcptf workspace create -org=my-org -json-input=@workspace.json

# From stdin
cat config.json | hcptf workspace create -org=my-org -json-input=-
```

Note: Routing flags like `-org` and `-workspace` are still required alongside `-json-input`.

### Field Filtering

Limit output to specific fields with `-fields` to reduce response size:

```bash
# Only return ID and Name
hcptf workspace list -org=my-org -fields=Name,ID -output=json
hcptf workspace read -org=my-org -name=prod -fields=ID,Name,Status
hcptf variable create -org=my-org -workspace=prod -key=x -value=y -fields=ID,Key
```

### Input Validation

All mutation commands validate inputs against common agent mistakes:
- Path traversal (`../`) is rejected
- Query injection (`?`, `&`) is rejected
- URL-encoded sequences (`%2f`) are rejected
- Control characters are rejected
- Overly long strings are rejected

### Recommended Agent Workflow

```bash
# 1. Discover flags for the command
hcptf schema workspace create

# 2. Dry-run to validate
hcptf workspace create -org=my-org -name=staging -dry-run

# 3. Execute with JSON output and field filtering
hcptf workspace create -org=my-org -name=staging -output=json -fields=ID,Name

# 4. For complex mutations, use JSON input
hcptf workspace create -org=my-org \
  -json-input='{"name":"staging","terraform-version":"1.6.0","auto-apply":false}' \
  -output=json -fields=ID,Name
```

## Output Formats

### Table Output (Default)

```bash
hcptf workspace list -org=my-org
```

### JSON Output (for scripting and agents)

```bash
hcptf workspace list -org=my-org -output=json
```

### Filtered JSON Output

```bash
hcptf workspace list -org=my-org -output=json -fields=Name,ID,Status
```

## Best Practices

1. **Always use `-dry-run` before write/delete commands** to confirm the action is correct before executing

2. **Always confirm with user before write/delete commands** — never auto-execute mutations without user approval

3. **Use `hcptf schema` to discover flags** — never guess flag names or parse help text

4. **Use `-output=json` and `-fields` for automation** — reduces output size and gives structured data

5. **Use hierarchical commands for clarity**
   - `hcptf workspace run create` for run workflows
   - `hcptf workspace run assessmentresult list` for drift workflows

6. **Prefer explicit flags for automation**
   - Use `-org=` and `-name=` in scripts
   - `-workspace=` remains available as alias
   - URL-style is great for interactive use

7. **Check status before applying**

```bash
# Always review before apply
hcptf workspace run show -id=run-abc123
hcptf workspace run apply -id=run-abc123
```

8. **Sensitive variables**

   ```bash
   # Always use -sensitive for secrets
   hcptf variable create -org=my-org -workspace=prod \
     -key=API_KEY -value=secret -sensitive
   ```

9. **Check drift before deployment**

```bash
# Assessment results show drift
hcptf workspace run assessmentresult list -org=my-org -name=prod
hcptf workspace run assessmentresult read -id=ar-abc123 -show-drift
hcptf workspace run create -org=my-org -name=prod -refresh-only
```

## Common Flags

| Flag            | Alias        | Description                          |
| --------------- | ------------ | ------------------------------------ |
| `-organization` | `-org`       | Organization name                    |
| `-name`         | `-workspace` | Workspace name                       |
| `-output`       |              | `table` (default) or `json`          |
| `-force`        |              | Skip confirmation prompts            |
| `-dry-run`      |              | Validate without making API calls    |
| `-fields`       |              | Comma-separated output field filter  |
| `-json-input`   |              | JSON payload (inline, `@file`, `-`)  |

## Getting Help

```bash
# General help
hcptf --help

# Command group help
hcptf registry --help
hcptf stack --help

# Specific command help
hcptf workspace create --help
hcptf registry module list --help
```

## Error Handling

When commands fail:

1. **Check authentication**: Ensure `TFE_TOKEN` is set or run `hcptf login`
2. **Verify organization**: Confirm you have access to the org
3. **Check resource names**: Workspaces and runs must exist
4. **Review permissions**: Ensure your token has required permissions

## Examples

### Complete Workspace Setup

```bash
# Create workspace with VCS
hcptf workspace create -org=my-org -name=api-service \
  -vcs-repo=owner/repo \
  -vcs-branch=main \
  -working-directory=terraform/

# Configure variables
hcptf variable create -org=my-org -workspace=api-service \
  -key=environment -value=staging

hcptf variable create -org=my-org -workspace=api-service \
  -key=AWS_ACCESS_KEY_ID -value=$AWS_KEY -category=env -sensitive

# Set up notifications
hcptf notification create -org=my-org -workspace=api-service \
  -name=slack-alerts \
  -destination-type=slack \
  -url=$SLACK_WEBHOOK \
  -triggers=run:completed,run:errored
```

### State Management

```bash
# List state versions
hcptf state list -org=my-org -workspace=production

# View state outputs
hcptf state outputs -org=my-org -workspace=production

# Export state
hcptf state read -org=my-org -workspace=production -output=json > state.json
```

### Policy Management

```bash
# Create policy
hcptf policy create -org=my-org \
  -name=cost-limit \
  -enforce-mode=hard-mandatory \
  -policy-file=policies/cost.sentinel

# Create policy set
hcptf policyset create -org=my-org \
  -name=production-policies \
  -global=true

# Add policy to set
hcptf policyset add-policy -org=my-org \
  -policyset=production-policies \
  -policy=cost-limit
```

## Troubleshooting

### Authentication Issues

```bash
# Verify token is valid
hcptf account show

# Re-login if needed
hcptf logout
hcptf login
```

### Finding Resource IDs

```bash
# Get workspace ID
hcptf workspace read -org=my-org -workspace=staging -output=json | \
  jq -r '.data.id'

# Get latest run ID
hcptf workspace run list -org=my-org -name=staging -output=json | \
  jq -r '.data[0].id'
```

### Debugging

```bash
# Use JSON output for detailed info
hcptf workspace run show -id=run-abc123 -output=json | jq .

# Check drift in assessment results
hcptf workspace run assessmentresult read -id=ar-abc123 -show-drift -summary-only
```

## Reference

- Authentication: See `hcptf login --help`
- Configuration: `~/.hcptfrc` (HCL format)
- Documentation: README.md in repository
- API Coverage: 100+ commands across 50+ resource types

## Notes

- All commands support `-help` flag for detailed usage
- Hierarchical structure (registry, stack) is the current standard
- Legacy flat commands have been removed in favor of hierarchical
- URL-style navigation works alongside traditional commands
