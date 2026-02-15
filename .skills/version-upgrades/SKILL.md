# Version Upgrade Skill

## Overview
This skill helps upgrade Terraform, provider, module, and policy versions in HCP Terraform workspaces. Keeping versions current improves security, performance, and access to new features.

## Prerequisites
- Authenticated with `hcptf` CLI
- Write access to target workspace
- VCS access if updating provider/module/policy versions in code

## Core Concepts

**Version Types:**
- **Terraform Version**: The Terraform CLI version used to run plans/applies
- **Provider Versions**: API plugins defined in `required_providers` block
- **Module Versions**: Reusable module versions called with `source` argument
- **Policy Versions**: Sentinel/OPA policy sets for compliance and governance

**Upgrade Approaches:**

1. **Terraform Version** (Simple - Workspace Setting):
   - One command: `hcptf workspace update -terraform-version=latest`
   - No code changes needed
   - Takes effect on next run

2. **Provider/Module/Policy Versions** (Complex - Requires Code Changes):
   - Get git repo info from workspace
   - Clone repository
   - Edit Terraform files (versions.tf, main.tf, sentinel.hcl, etc.)
   - Update version constraints
   - Commit and push changes
   - VCS triggers new run automatically

## Workflow

### 1. Check Current Versions

**For a specific workspace:**

```bash
# View workspace Terraform version
hcptf <org> <workspace>
# Shows: TerraformVersion: ~>1.9.0

# Get detailed version info via Explorer
hcptf explorer query -org=<org> -type=workspaces \
  -filter="workspace-name:<workspace>" \
  -fields=workspace-name,workspace-terraform-version,providers,modules
```

**For all workspaces in org:**

```bash
# Find outdated Terraform versions
hcptf explorer query -org=<org> -type=tf_versions \
  -sort=version

# Find outdated providers
hcptf explorer query -org=<org> -type=providers \
  -sort=-version \
  -fields=name,version,workspace-count,workspaces

# Find outdated modules
hcptf explorer query -org=<org> -type=modules \
  -sort=-version \
  -fields=name,version,workspace-count,workspaces
```

### 2. Find Latest Versions

**Terraform versions:**
- Official releases: https://releases.hashicorp.com/terraform/
- What's in use across org: `hcptf explorer query -org=<org> -type=tf_versions`

**Public provider versions:**
- Provider registry page: `https://registry.terraform.io/providers/<namespace>/<name>/latest`
  - Example: https://registry.terraform.io/providers/hashicorp/aws/latest
  - Example: https://registry.terraform.io/providers/hashicorp/random/latest
- Documentation includes changelog, upgrade guides, and version history
- In use across org: `hcptf explorer query -org=<org> -type=providers`
- CLI commands:
  ```bash
  # Get latest provider version and details
  hcptf publicregistry provider -name=hashicorp/aws

  # List all available provider versions
  hcptf publicregistry provider versions -name=hashicorp/aws
  ```

**Private provider versions:**
```bash
# List private providers
hcptf registry provider list -organization=<org>

# View provider details and versions
hcptf registry provider read -organization=<org> -name=<provider>

# View specific version
hcptf registry provider version read -organization=<org> \
  -name=<provider> -version=<version>
```

**Public module versions:**
- Module registry page: `https://registry.terraform.io/modules/<namespace>/<name>/<system>`
  - Example: https://registry.terraform.io/modules/terraform-aws-modules/s3-bucket/aws
  - Example: https://registry.terraform.io/modules/hashicorp/dir/template
- Shows available versions, inputs/outputs, and usage examples
- In use across org: `hcptf explorer query -org=<org> -type=modules`
- CLI commands:
  ```bash
  # Get latest module version and details
  hcptf publicregistry module -name=terraform-aws-modules/vpc/aws

  # Check module downloads, verified status, docs URL
  hcptf publicregistry module -name=terraform-aws-modules/s3-bucket/aws
  ```

**Private module versions:**
```bash
# List private modules
hcptf registry module list -organization=<org>

# View module details and versions
hcptf registry module read -organization=<org> -namespace=<org> -name=<module>
```

**Public policy versions:**
- Policy registry page: `https://registry.terraform.io/policies/<namespace>/<name>`
  - Example: https://registry.terraform.io/policies/hashicorp/CIS-Policy-Set-for-AWS-Terraform
  - Example: https://registry.terraform.io/policies/hashicorp/gcp-networking-terraform
- Shows available versions, included policies, and modules
- CLI commands:
  ```bash
  # List all available public policies
  hcptf publicregistry policy list

  # Get latest policy version and details
  hcptf publicregistry policy -name=hashicorp/CIS-Policy-Set-for-AWS-Terraform

  # Get specific policy version
  hcptf publicregistry policy -name=hashicorp/CIS-Policy-Set-for-Azure-Terraform -version=1.0.0

  # Check policy count, module count, included policies
  hcptf publicregistry policy -name=hashicorp/gcp-networking-terraform
  ```

**Private policy versions:**
```bash
# List private policies (via policy sets)
hcptf policyset list -organization=<org>

# View policy set details
hcptf policyset read -organization=<org> -id=<policy-set-id>
```

### 3. Upgrade Terraform Version (Workspace Setting)

```bash
# Update workspace to use latest Terraform
hcptf workspace update -org=<org> -name=<workspace> \
  -terraform-version=1.10.0

# Or use "latest" to always use newest
hcptf workspace update -org=<org> -name=<workspace> \
  -terraform-version=latest

# Or use version constraint
hcptf workspace update -org=<org> -name=<workspace> \
  -terraform-version="~>1.10.0"
```

### 4. Upgrade Provider/Module/Policy Versions (Code Changes)

Provider, module, and policy versions are defined in code and require updating the Terraform files.

**Get code location:**

```bash
# Get current configuration
RUN_ID=$(hcptf <org> <workspace> -output=json | jq -r '.CurrentRunID')
hcptf <org> <workspace> runs $RUN_ID configversion
# Shows: RepoIdentifier, Branch, CommitURL, CommitSHA
```

**Update provider version in code:**

```hcl
# In versions.tf or terraform block
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.28.0"  # Updated from 5.69.0
    }
  }
}
```

**Update module version in code:**

```hcl
# In main.tf or module call
module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> 5.10.0"  # Updated from 4.0.0
  # ...
}
```

### 5. Test the Upgrade

Create a test run to verify the upgrade works:

```bash
# Create a run with the updated versions
hcptf <org> <workspace> runs create \
  -message="Upgrade Terraform to 1.10.0, AWS provider to 6.28.0"

# Monitor the run
hcptf <org> <workspace> runs list | head -5

# Check for errors
hcptf <org> <workspace> runs <run-id> show
```

If the plan shows issues:
- Check plan output: `hcptf <org> <workspace> runs <run-id> plan`
- View logs: `hcptf <org> <workspace> runs <run-id> logs`

### 6. Apply if Successful

```bash
# If plan looks good, apply
hcptf <org> <workspace> runs <run-id> apply

# Or if auto-apply is enabled, run is applied automatically
```

### 7. Verify

```bash
# Check updated versions
hcptf <org> <workspace>
# Shows: TerraformVersion: 1.10.0

hcptf explorer query -org=<org> -type=workspaces \
  -filter="workspace-name:<workspace>" \
  -fields=workspace-name,workspace-terraform-version,providers,modules
```

## Common Scenarios

### Scenario 1: Upgrade Terraform Version Only

```bash
# 1. Check current version
hcptf my-org my-workspace
# Shows: TerraformVersion: 1.9.6

# 2. Upgrade to latest
hcptf workspace update -org=my-org -name=my-workspace \
  -terraform-version=latest

# 3. Test with a plan (don't apply yet)
hcptf my-org my-workspace runs create -message="Test Terraform 1.10 upgrade"

# 4. Verify plan succeeds, then apply
hcptf my-org my-workspace runs <run-id> apply
```

### Scenario 2: Upgrade AWS Provider (Major Version)

**Important**: Provider versions are in code, not workspace settings. You MUST clone the repo and edit files.

```bash
# 1. Check current provider version
hcptf explorer query -org=my-org -type=workspaces \
  -filter="workspace-name:my-workspace" \
  -fields=workspace-name,providers
# Shows: hashicorp/aws:5.69.0

# 2. Get git repo information
RUN_ID=$(hcptf my-org my-workspace -output=json | jq -r '.CurrentRunID')
hcptf my-org my-workspace runs $RUN_ID configversion
# Output:
#   RepoIdentifier: thrashr888/my-infra
#   Branch: main
#   CommitSHA: abc123...
#   CommitURL: https://github.com/thrashr888/my-infra/commit/abc123...

# 3. Clone and update code
git clone https://github.com/thrashr888/my-infra
cd my-infra
git checkout main

# Review upgrade guide BEFORE editing (check for breaking changes):
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-6-upgrade
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-6-changelog

# Edit versions.tf or terraform block:
# Change:
#   aws = { version = "~> 5.69.0" }
# To:
#   aws = { version = "~> 6.0" }

# May also need to update resource configurations if provider API changed

# 4. Commit and push
git add versions.tf
git commit -m "Upgrade AWS provider from v5 to v6"
git push origin main
# VCS push triggers run automatically in HCP Terraform

# 5. Monitor the automatically triggered run
hcptf my-org my-workspace runs list | head -5
RUN_ID=$(hcptf my-org my-workspace -output=json | jq -r '.CurrentRunID')

# 6. Review plan for provider upgrade impacts
hcptf my-org my-workspace runs $RUN_ID show

# 7. Apply if successful (or it auto-applies if enabled)
hcptf my-org my-workspace runs $RUN_ID apply

# 8. Verify new provider version
hcptf explorer query -org=my-org -type=workspaces \
  -filter="workspace-name:my-workspace" \
  -fields=workspace-name,providers
# Should now show: hashicorp/aws:6.x.x
```

### Scenario 3: Upgrade Multiple Workspaces

```bash
# 1. Find all workspaces using old Terraform version
hcptf explorer query -org=my-org -type=tf_versions \
  -fields=version,workspace-count,workspaces
# Identify workspaces on 1.9.x

# 2. Bulk update Terraform version
for ws in workspace1 workspace2 workspace3; do
  echo "Upgrading $ws..."
  hcptf workspace update -org=my-org -name=$ws -terraform-version=1.10.0
done

# 3. Test each workspace
for ws in workspace1 workspace2 workspace3; do
  echo "Testing $ws..."
  hcptf my-org $ws runs create -message="Test Terraform 1.10 upgrade"
done
```

### Scenario 4: Upgrade Module Versions

**Important**: Module versions are in code. You MUST clone the repo and edit the module source blocks.

```bash
# 1. Find modules in use
hcptf explorer query -org=my-org -type=modules \
  -fields=name,version,workspaces
# Shows: terraform-aws-modules/s3-bucket/aws version 4.0.0 used in cool-website

# 2. Get git repo for workspace using outdated module
RUN_ID=$(hcptf my-org cool-website -output=json | jq -r '.CurrentRunID')
hcptf my-org cool-website runs $RUN_ID configversion
# Output:
#   RepoIdentifier: thrashr888/cool-website
#   Branch: main
#   CommitURL: https://github.com/thrashr888/cool-website/commit/abc123...

# 3. Clone and update code
git clone https://github.com/thrashr888/cool-website
cd cool-website
git checkout main

# Check module docs BEFORE editing (for breaking changes):
# https://registry.terraform.io/modules/terraform-aws-modules/s3-bucket/aws/latest
# Review: Inputs, Outputs, Changelog, Examples

# Edit main.tf (or wherever module is called):
# Change:
#   module "s3_bucket" {
#     source  = "terraform-aws-modules/s3-bucket/aws"
#     version = "~> 4.0"
#     bucket  = var.bucket_name
#     ...
#   }
# To:
#   module "s3_bucket" {
#     source  = "terraform-aws-modules/s3-bucket/aws"
#     version = "~> 5.10"
#     bucket  = var.bucket_name
#     # Add any new required variables
#     # Remove any deprecated variables
#     ...
#   }

# 4. Commit and push
git add main.tf
git commit -m "Upgrade s3-bucket module from v4 to v5.10"
git push origin main
# VCS push triggers run automatically

# 5. Monitor and verify
hcptf my-org cool-website runs list | head -5
hcptf explorer query -org=my-org -type=modules \
  -fields=name,version,workspaces | grep cool-website
```

### Scenario 5: Upgrade Policy Versions

**Important**: Policy versions are managed in VCS-backed policy sets via sentinel.hcl configuration.

```bash
# 1. Find current policy sets in use
hcptf policyset list -organization=my-org

# 2. Check for newer policy versions
hcptf publicregistry policy list | grep CIS

# 3. Get details about latest policy version
hcptf publicregistry policy -name=hashicorp/CIS-Policy-Set-for-AWS-Terraform
# Shows: Version: 1.0.1, PolicyCount: 35, ModuleCount: 4

# 4. Review policy changes
# Visit: https://registry.terraform.io/policies/hashicorp/CIS-Policy-Set-for-AWS-Terraform
# Review: Changelog, new policies, removed policies, parameter changes

# 5. Get VCS info for policy set repository
hcptf policyset read -organization=my-org -id=<policy-set-id>
# Shows VCS repo and branch information

# 6. Clone and update policy configuration
git clone <policy-repo-url>
cd policy-repo

# Edit sentinel.hcl to update policy source version:
# Change:
#   policy "require-mfa" {
#     source = "https://registry.terraform.io/v2/policies/hashicorp/CIS-Policy-Set-for-AWS-Terraform/1.0.0/policy-modules/..."
#   }
# To:
#   policy "require-mfa" {
#     source = "https://registry.terraform.io/v2/policies/hashicorp/CIS-Policy-Set-for-AWS-Terraform/1.0.1/policy-modules/..."
#   }

# 7. Commit and push
git add sentinel.hcl
git commit -m "Upgrade CIS AWS policy set from v1.0.0 to v1.0.1"
git push origin main
# Policy set updates automatically in HCP Terraform

# 8. Verify policy set updated
hcptf policyset read -organization=my-org -id=<policy-set-id>
```

## Version Compatibility

**Terraform version constraints:**
- `1.9.0` - Exact version
- `~> 1.9.0` - Pessimistic constraint (1.9.x only)
- `~> 1.9` - Allow 1.x (1.9.0, 1.10.0, etc.)
- `>= 1.9.0` - Minimum version
- `latest` - Always use newest available

**Provider version constraints:**
Follow same syntax. Use pessimistic constraints to avoid breaking changes:
- `~> 5.69.0` - Allow 5.69.x patches only
- `~> 5.0` - Allow 5.x minor/patch updates

## Upgrade Considerations

**Before upgrading:**
1. Review changelogs for breaking changes
2. Check workspace health (no drift, checks passing)
3. Have rollback plan ready
4. Test in non-production workspace first

**Terraform version upgrades:**
- Usually safe for patch versions (1.9.6 → 1.9.7)
- Minor versions may have deprecations (1.9.x → 1.10.x)
- Major versions have breaking changes (review upgrade guide)

**Provider version upgrades:**
- Patch versions are safe (5.69.0 → 5.69.1)
- Minor versions may add deprecations (5.69.x → 5.70.x)
- Major versions have breaking changes (5.x → 6.x)
- Always review provider upgrade guide

**Module version upgrades:**
- Depends on module's versioning practice
- Check module changelog and README
- Test for interface changes (new required variables, removed outputs)

**Policy version upgrades:**
- Review policy changelog for new/removed checks
- Check if new policies require additional parameters
- Test in non-production workspace first
- Consider impact on existing runs (advisory vs mandatory enforcement)
- Verify policy modules are compatible

## Rollback

**If upgrade fails:**

```bash
# Revert Terraform version
hcptf workspace update -org=<org> -name=<workspace> \
  -terraform-version=1.9.6

# Revert code changes
#   - Git revert the commit
#   - Push to trigger new run
#   - Or discard the failed run and create new one with old config
```

## Related Commands

- `hcptf explorer query` - Find versions in use across org
- `hcptf workspace read` - View workspace Terraform version
- `hcptf workspace update` - Update workspace Terraform version
- `hcptf configversion read` - Get VCS info to find code
- `hcptf run create` - Test upgraded configuration
- `hcptf run show` - Monitor upgrade run
- `hcptf publicregistry provider` - Get public provider info and latest version
- `hcptf publicregistry provider versions` - List all available provider versions
- `hcptf publicregistry module` - Get public module info and latest version
- `hcptf publicregistry policy` - Get public policy info and latest version
- `hcptf publicregistry policy list` - List all available public policies
- `hcptf registry module read` - Check private module versions
- `hcptf policyset list` - List policy sets in organization
- `hcptf policyset read` - View policy set details and VCS info

## Tips

1. **Upgrade incrementally**: Don't jump multiple major versions at once
2. **Test first**: Always run plan before apply
3. **Document**: Include clear messages in runs explaining what's being upgraded
4. **Monitor**: Watch run logs for deprecation warnings
5. **Coordinate**: For major provider upgrades, update all workspaces using that provider
6. **Use constraints**: Prefer `~>` constraints over exact versions for easier patch updates
7. **Stage upgrades**: Test in dev/staging before production
8. **Policy upgrades**: Review new policy checks and set to advisory mode first before enforcing
9. **Use CLI tools**: Leverage `publicregistry` commands to quickly check latest versions without leaving terminal
