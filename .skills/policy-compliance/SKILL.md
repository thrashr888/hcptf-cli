---
name: policy-compliance
description: Investigate and resolve policy check failures in HCP Terraform. Use when runs are blocked by failed Sentinel or OPA policy checks, understanding why policies failed, deciding whether to fix code or override policies, or finding policy failures across the organization.
---

# Policy Compliance Skill

## Overview

This skill helps investigate and resolve policy check failures in HCP Terraform. Policies enforce governance, security, and compliance standards using Sentinel or OPA. When policies fail, you need to understand what was checked, why it failed, and how to remediate.

## Prerequisites

- Authenticated with `hcptf` CLI
- Read access to workspaces and policy sets
- Override permissions if you need to bypass failed policies
- Write access to workspace code if fixing violations

## Core Concepts

**Policy Types:**

- **Sentinel**: HashiCorp's policy-as-code framework
- **OPA (Open Policy Agent)**: CNCF policy framework

**Enforcement Levels:**

- **Advisory**: Logs a warning, run continues (informational)
- **Soft-Mandatory**: Can be overridden by authorized users
- **Hard-Mandatory**: Cannot be overridden, run fails

**Policy Check Statuses:**

- **Passed**: All policies passed
- **Failed**: One or more mandatory policies failed
- **Overridden**: Failed policy was manually overridden
- **Unreachable**: Policy could not be evaluated (configuration error)

## Workflow

### 1. Detect Policy Failures

**From run status:**

```bash
# Check run status (shows policy check summary)
hcptf <org> <workspace> runs <run-id> show

# Look for:
#   PolicyChecks: failed
#   Status: policy_checked (awaiting override)
```

**Across organization:**

```bash
# Find all runs with policy failures
hcptf queryrun list -organization=<org> -status=policy_checked

# Find workspaces with recent policy failures
hcptf explorer query -org=<org> -type=workspaces \
  -fields=workspace-name,current-run-status
```

### 2. View Policy Check Details

```bash
# List all policy checks for a run
hcptf policycheck list -run-id=<run-id>

# Shows:
#   - Policy check ID
#   - Status (passed/failed/overridden)
#   - Total policies passed/failed
#   - Enforcement level

# Get detailed results for a specific policy check
hcptf policycheck read -id=<policy-check-id>

# Shows:
#   - Individual policy results
#   - Failed policy names
#   - Policy evaluation results
#   - Which resources violated policies
```

**URL-style navigation:**

```bash
# View policy checks for a run
hcptf <org> <workspace> runs <run-id> policychecks
```

### 3. Understand What the Policy Does

**For public policies:**

```bash
# Search for the policy
hcptf publicregistry policy list | grep CIS

# Get policy details
hcptf publicregistry policy -name=hashicorp/CIS-Policy-Set-for-AWS-Terraform

# Shows:
#   - PolicyCount: 35
#   - Policies: cloudtrail-logs-bucket-not-public, efs-encryption-at-rest-enabled, ...
#   - DocsURL: https://registry.terraform.io/policies/...

# Visit docs URL to see:
#   - What each policy checks
#   - Why it matters (security/compliance)
#   - How to fix violations
#   - Parameters and configuration
```

**For private/custom policies:**

```bash
# Get policy set details
hcptf policyset read -organization=<org> -id=<policy-set-id>

# Shows:
#   - VCS repository with policy code
#   - Policy configuration
#   - Enforcement levels
#   - Workspaces where applied
```

### 4. Identify Violating Resources

Policy check results show which Terraform resources violated policies:

```bash
# View policy check details
hcptf policycheck read -id=<policy-check-id>

# Look for output like:
#   Policy: ec2-security-group-ingress-traffic-restriction-port-22
#   Status: failed
#   Message: Security group allows SSH from 0.0.0.0/0
#   Resource: aws_security_group.web
```

**Cross-reference with plan:**

```bash
# View the plan to see resource configuration
hcptf <org> <workspace> runs <run-id> plan

# Find the violating resource and see what's being created/changed
```

### 5. Decide on Remediation Strategy

**Decision Matrix:**

| Situation                          | Action                      | Command                                            |
| ---------------------------------- | --------------------------- | -------------------------------------------------- |
| Policy is correct, code is wrong   | Fix code                    | Clone repo, edit .tf files, commit                 |
| Policy is too strict for this case | Override policy             | `hcptf policycheck override`                       |
| Policy needs adjustment            | Update policy set           | Edit policy in VCS, update parameters              |
| Need temporary exception           | Override with justification | `hcptf policycheck override` (requires permission) |
| Policy is informational            | Acknowledge and continue    | Review advisory warnings, no action needed         |

### 6. Remediation Options

**Option A: Fix the Code**

Most common approach - update Terraform code to comply with policy.

```bash
# 1. Get VCS info for workspace
RUN_ID=$(hcptf <org> <workspace> -output=json | jq -r '.CurrentRunID')
hcptf <org> <workspace> runs $RUN_ID configversion

# Shows:
#   RepoIdentifier: my-org/infrastructure
#   Branch: main
#   CommitSHA: abc123...

# 2. Clone repository
git clone https://github.com/my-org/infrastructure
cd infrastructure

# 3. Fix the violation
# Example: Security group allowing SSH from 0.0.0.0/0
# Edit security_groups.tf:
# Change:
#   ingress {
#     from_port   = 22
#     to_port     = 22
#     protocol    = "tcp"
#     cidr_blocks = ["0.0.0.0/0"]  # ❌ Too permissive
#   }
# To:
#   ingress {
#     from_port   = 22
#     to_port     = 22
#     protocol    = "tcp"
#     cidr_blocks = ["10.0.0.0/8"]  # ✅ Restrict to private network
#   }

# 4. Commit and push
git add security_groups.tf
git commit -m "Restrict SSH access to private network for policy compliance"
git push origin main

# VCS push triggers new run automatically
# 5. Verify new run passes policy checks
hcptf <org> <workspace> runs list | head -5
hcptf <org> <workspace> runs <new-run-id> show
```

**Option B: Override the Policy**

Use when policy is too strict for a specific case. Requires override permissions.

```bash
# Override a failed policy check
hcptf policycheck override -id=<policy-check-id>

# Run continues and can be applied
hcptf <org> <workspace> runs <run-id> apply

# Note: Only works for soft-mandatory policies
# Hard-mandatory policies cannot be overridden
```

**Option C: Adjust Policy Enforcement**

Update the policy set to change enforcement level or parameters.

```bash
# 1. Get policy set info
hcptf policyset read -organization=<org> -id=<policy-set-id>

# 2. Clone policy repository
git clone <policy-vcs-repo>
cd policy-repo

# 3. Edit sentinel.hcl to adjust enforcement
# Change:
#   policy "require-tags" {
#     enforcement_level = "hard-mandatory"
#   }
# To:
#   policy "require-tags" {
#     enforcement_level = "soft-mandatory"  # Allow overrides
#   }

# Or add parameters to make policy less strict:
#   params = {
#     allowed_regions = ["us-east-1", "us-west-2", "eu-west-1"]
#   }

# 4. Commit and push
git add sentinel.hcl
git commit -m "Allow overrides for tagging policy in dev workspaces"
git push origin main

# Policy set updates automatically
# 5. Trigger a new run to re-evaluate with updated policy
hcptf <org> <workspace> runs create -message="Re-run with updated policy"
```

## Common Policy Scenarios

### Scenario 1: CIS Benchmark Violation

```bash
# Run fails with CIS AWS policy
hcptf <org> production-app runs run-abc123 policychecks

# Shows failure:
#   Policy: cloudtrail-bucket-access-logging-enabled
#   Status: failed
#   Enforcement: hard-mandatory

# Understand the policy
hcptf publicregistry policy -name=hashicorp/CIS-Policy-Set-for-AWS-Terraform
# Visit DocsURL to see CIS 3.6 requirement

# Check which resource failed
hcptf policycheck read -id=polchk-xyz789
# Output: aws_s3_bucket.cloudtrail_logs missing logging configuration

# Fix the code
# Add to terraform:
resource "aws_s3_bucket_logging" "cloudtrail_logs" {
  bucket = aws_s3_bucket.cloudtrail_logs.id

  target_bucket = aws_s3_bucket.log_bucket.id
  target_prefix = "cloudtrail/"
}

# Commit, push, verify new run passes
```

### Scenario 2: Tagging Policy Failure

```bash
# Advisory policy warns about missing tags
hcptf <org> dev-environment runs run-def456 policychecks

# Shows:
#   Policy: require-common-tags
#   Status: failed
#   Enforcement: advisory (run can continue)

# Check which resources are missing tags
hcptf policycheck read -id=polchk-abc123
# Output:
#   aws_instance.web - missing tags: Environment, Owner, CostCenter
#   aws_s3_bucket.data - missing tags: Environment, CostCenter

# Fix by adding tags
# In main.tf:
resource "aws_instance" "web" {
  # ... existing config ...

  tags = {
    Name        = "web-server"
    Environment = "development"
    Owner       = "platform-team"
    CostCenter  = "engineering"
  }
}

# Or use default_tags in provider:
provider "aws" {
  default_tags {
    tags = {
      Environment = "development"
      ManagedBy   = "Terraform"
      CostCenter  = "engineering"
    }
  }
}
```

### Scenario 3: Security Group Overly Permissive

```bash
# Hard-mandatory security policy fails
hcptf <org> api-service runs run-ghi789 show
# Status: policy_checked (blocked)

hcptf <org> api-service runs run-ghi789 policychecks

# Shows:
#   Policy: ec2-security-group-ingress-traffic-restriction-port-22
#   Status: failed
#   Enforcement: hard-mandatory (cannot override)

# Check details
hcptf policycheck read -id=polchk-def456
# Output: aws_security_group.api allows 0.0.0.0/0 on port 22

# MUST fix code (cannot override hard-mandatory):
# Change security_groups.tf:
resource "aws_security_group" "api" {
  name = "api-sg"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/8"]  # Restrict to VPN/private network
    description = "SSH from VPN"
  }

  # Keep application port public
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS from internet"
  }
}
```

### Scenario 4: Find All Policy Failures Across Organization

```bash
# Find workspaces with policy failures
hcptf explorer query -org=<org> -type=workspaces \
  -fields=workspace-name,current-run-status,current-run-id \
  -filter="current-run-status:policy_checked"

# Check each workspace's policy failures
for ws in workspace1 workspace2 workspace3; do
  echo "=== $ws ==="
  RUN_ID=$(hcptf <org> $ws -output=json | jq -r '.CurrentRunID')
  hcptf policycheck list -run-id=$RUN_ID
done

# Generate compliance report
hcptf explorer query -org=<org> -type=workspaces \
  -fields=workspace-name,current-run-status,policy-set-names \
  -output=json > compliance-report.json
```

### Scenario 5: Test Policy Changes Before Enforcement

```bash
# 1. Set policy to advisory mode first
# Edit sentinel.hcl:
policy "new-security-check" {
  enforcement_level = "advisory"  # Test first
}

# 2. Apply to test workspace
hcptf policyset update -organization=<org> -id=<policy-set-id> \
  -workspace=test-workspace

# 3. Trigger run and check results
hcptf test-workspace runs create -message="Test new security policy"

# 4. Review advisory failures (won't block run)
hcptf test-workspace runs <run-id> policychecks

# 5. Once validated, increase enforcement
# Edit sentinel.hcl:
policy "new-security-check" {
  enforcement_level = "soft-mandatory"  # Enforce with overrides
}

# 6. Later, make it hard-mandatory if needed
```

## Policy Troubleshooting

**Policy check shows "unreachable":**

- Policy code has syntax errors
- Policy cannot access required data (tfplan, tfconfig)
- Check policy set VCS repository for errors

**Override button not available:**

- Policy is hard-mandatory (cannot override)
- User lacks override permissions
- Check team access settings

**Policy passes locally but fails in HCP Terraform:**

- Sentinel CLI vs cloud environment differences
- Check policy parameters and configuration
- Verify policy set is correctly applied

**Policy fails on every run:**

- Policy configuration may be incorrect
- Check if policy matches workspace type (EC2 policy on GCP workspace)
- Review policy parameters

## Best Practices

1. **Understand Before Overriding**: Always investigate why a policy failed before overriding
2. **Document Overrides**: Include justification when overriding policies
3. **Advisory First**: Test new policies in advisory mode before enforcing
4. **Fix Code, Not Policy**: Prefer fixing violations over relaxing policies
5. **Regular Reviews**: Periodically review policy overrides and adjust policies if needed
6. **Workspace Exceptions**: Use different policy sets for dev vs production
7. **Incremental Rollout**: Apply new policies to test workspaces first
8. **Monitor Compliance**: Track policy failures across organization
9. **Educate Teams**: Share policy documentation and remediation guides

## Related Commands

- `hcptf policycheck list` - List policy checks for a run
- `hcptf policycheck read` - View detailed policy check results
- `hcptf policycheck override` - Override a failed policy (if allowed)
- `hcptf policyset list` - List policy sets in organization
- `hcptf policyset read` - View policy set configuration
- `hcptf publicregistry policy` - Get public policy details and documentation
- `hcptf publicregistry policy list` - Browse available public policies
- `hcptf explorer query` - Find workspaces with policy failures
- `hcptf queryrun list` - Find runs with specific statuses
- `hcptf <org> <workspace> runs <run-id> policychecks` - View policy checks (URL-style)

## Policy Compliance Metrics

Track compliance across your organization:

```bash
# Total workspaces with policy failures
hcptf explorer query -org=<org> -type=workspaces \
  -filter="current-run-status:policy_checked" \
  -fields=workspace-name | wc -l

# Policy sets in use
hcptf policyset list -organization=<org>

# Recent runs requiring policy override
hcptf queryrun list -organization=<org> -status=policy_override

# Workspaces by policy set
hcptf explorer query -org=<org> -type=workspaces \
  -fields=workspace-name,policy-set-names
```
