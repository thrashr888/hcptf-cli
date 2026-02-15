# Version Upgrade Skill

## Overview
This skill helps upgrade Terraform, provider, and module versions in HCP Terraform workspaces. Keeping versions current improves security, performance, and access to new features.

## Prerequisites
- Authenticated with `hcptf` CLI
- Write access to target workspace
- VCS access if updating provider/module versions in code

## Core Concepts

**Version Types:**
- **Terraform Version**: The Terraform CLI version used to run plans/applies (workspace setting)
- **Provider Versions**: API plugins defined in `required_providers` block (code)
- **Module Versions**: Reusable module versions called with `source` argument (code)

**Upgrade Approaches:**
- **Terraform Version**: Update workspace setting (no code change needed)
- **Providers/Modules**: Update version constraints in code, commit, and test

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
- Check https://releases.hashicorp.com/terraform/
- Or see what's in use across your org (Explorer tf_versions)

**Provider versions:**
- Check https://registry.terraform.io/providers/hashicorp/aws/latest
- Or see highest version in use across org (Explorer providers)

**Module versions:**
- Check public registry: https://registry.terraform.io/modules/<namespace>/<name>/<system>
- Or check private registry: `hcptf registry module read`

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

### 4. Upgrade Provider/Module Versions (Code Changes)

Provider and module versions are defined in code and require updating the Terraform files.

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

```bash
# 1. Check current provider version
hcptf explorer query -org=my-org -type=workspaces \
  -filter="workspace-name:my-workspace" \
  -fields=workspace-name,providers
# Shows: hashicorp/aws:5.69.0

# 2. Get code location
RUN_ID=$(hcptf my-org my-workspace -output=json | jq -r '.CurrentRunID')
hcptf my-org my-workspace runs $RUN_ID configversion
# Shows: RepoIdentifier: user/repo, Branch: main, CommitURL: ...

# 3. Update code
#    - Clone repo: git clone https://github.com/user/repo
#    - Checkout branch: git checkout main
#    - Edit versions.tf: change aws version from "~> 5.69.0" to "~> 6.0"
#    - Review AWS provider upgrade guide for breaking changes
#    - Commit and push

# 4. Test upgrade
#    VCS push triggers run automatically, or:
hcptf my-org my-workspace runs create -message="Upgrade AWS provider to v6"

# 5. Review plan for provider upgrade impacts
hcptf my-org my-workspace runs <run-id> plan

# 6. Apply if successful
hcptf my-org my-workspace runs <run-id> apply
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

```bash
# 1. Find modules in use
hcptf explorer query -org=my-org -type=modules \
  -fields=name,version,workspaces

# 2. For each workspace using outdated module:
#    - Get code location
#    - Clone repo
#    - Update module version in .tf files
#    - Check module changelog for breaking changes
#    - Update module arguments if needed
#    - Commit and push

# 3. Test
hcptf my-org my-workspace runs create \
  -message="Upgrade s3-bucket module to 5.10.0"
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
- `hcptf registry module read` - Check private module versions

## Tips

1. **Upgrade incrementally**: Don't jump multiple major versions at once
2. **Test first**: Always run plan before apply
3. **Document**: Include clear messages in runs explaining what's being upgraded
4. **Monitor**: Watch run logs for deprecation warnings
5. **Coordinate**: For major provider upgrades, update all workspaces using that provider
6. **Use constraints**: Prefer `~>` constraints over exact versions for easier patch updates
7. **Stage upgrades**: Test in dev/staging before production
