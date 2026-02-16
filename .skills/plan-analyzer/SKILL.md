---
name: plan-analyzer
description: Analyze Terraform plan results before applying changes. Use when reviewing plans for safety, understanding impact, identifying risks, checking for unintended changes, or validating infrastructure modifications in HCP Terraform workspaces.
---

# Plan Analyzer Skill

## Overview

This skill helps review Terraform plan results to understand what will change, identify potential risks, and validate that changes match expectations before applying. It focuses on analyzing plan output to catch unintended changes, destructive operations, and configuration issues.

## Prerequisites

- HCP Terraform account with workspace access
- Authenticated with `hcptf` CLI (`TFE_TOKEN` or `~/.terraform.d/credentials.terraform.io`)
- Read access to workspace runs
- Active or completed plans to analyze

## Core Concepts

**Plan Analysis Areas:**

1. **Change Summary**: Understand the scope (additions, changes, destructions)
2. **Destructive Changes**: Identify resources that will be deleted or replaced
3. **High-Risk Changes**: Database deletions, security group changes, network reconfigurations
4. **Unintended Changes**: Catch configuration drift or unexpected modifications
5. **Resource Dependencies**: Understand change propagation and cascading effects
6. **Configuration Validation**: Verify values, names, and settings match expectations

**Common Issues to Check:**

- Accidental resource deletions (especially stateful resources like databases)
- Security group changes that open ports unintentionally
- Changes to production resources that should be stable
- Tag or naming changes that affect organization or billing
- Resource replacements (delete + create) that cause downtime
- Large-scale changes from small configuration tweaks
- Sensitive value changes (passwords, keys, certificates)

## Workflow

### 1. Find Runs with Plans

```bash
# List recent runs for a workspace
hcptf run list -org=my-org -workspace=my-workspace

# Or using URL-style
hcptf my-org my-workspace runs

# Find runs that are in planned state
hcptf run list -org=my-org -workspace=my-workspace -output=json | \
  jq -r '.[] | select(.Status == "planned") | "\(.ID) \(.CreatedAt)"'
```

### 2. Get Plan Summary

**Flag-based commands:**

```bash
# Read plan details by run ID
hcptf plan read -id=run-abc123

# Or by plan ID
hcptf plan read -id=plan-xyz789

# JSON output for parsing
hcptf plan read -id=run-abc123 -output=json | jq '{
  status: .Status,
  has_changes: .HasChanges,
  additions: .ResourceAdditions,
  changes: .ResourceChanges,
  destructions: .ResourceDestructions,
  imports: .ResourceImports
}'
```

**URL-style commands (alternative):**

```bash
# Plan summary via URL-style navigation
hcptf my-org my-workspace runs run-abc123 plan

# With JSON output
hcptf my-org my-workspace runs run-abc123 plan -output=json
```

### 3. Review Plan Text Output

**Flag-based commands:**

```bash
# Get full plan logs
hcptf plan logs -id=plan-abc123

# Save to file for detailed review
hcptf plan logs -id=run-abc123 > plan-output.txt
```

**URL-style commands (alternative):**

```bash
# Get plan logs via URL-style (shorter!)
hcptf my-org my-workspace runs run-abc123 logs

# Save to file
hcptf my-org my-workspace runs run-abc123 logs > plan-output.txt
```

**Searching plan output:**

```bash
# Search for specific resource types
hcptf plan logs -id=run-abc123 | grep "aws_db_instance"

# Find all resources being destroyed
hcptf plan logs -id=run-abc123 | grep -A 5 "will be destroyed"

# Find all resources being replaced
hcptf plan logs -id=run-abc123 | grep -A 5 "must be replaced"

# Check for security group changes
hcptf plan logs -id=run-abc123 | grep -A 10 "aws_security_group"
```

### 4. Analyze Change Patterns

**Check for destructive changes:**

```bash
# Get plan stats
PLAN_JSON=$(hcptf plan read -id=run-abc123 -output=json)

DESTRUCTIONS=$(echo "$PLAN_JSON" | jq -r '.ResourceDestructions')

if [ "$DESTRUCTIONS" -gt 0 ]; then
  echo "WARNING: $DESTRUCTIONS resources will be destroyed"
  echo "Review carefully:"
  hcptf plan logs -id=run-abc123 | grep -B 2 "will be destroyed"
fi
```

**Check for replacements (risky):**

```bash
# Replacements = resources that will be deleted then recreated
hcptf plan logs -id=run-abc123 | grep "must be replaced" | wc -l

# Show what's being replaced and why
hcptf plan logs -id=run-abc123 | grep -B 5 -A 10 "must be replaced"
```

**Check resource change ratio:**

```bash
# Large number of changes might indicate a problem
PLAN_JSON=$(hcptf plan read -id=run-abc123 -output=json)

TOTAL=$(echo "$PLAN_JSON" | jq -r '(.ResourceAdditions + .ResourceChanges + .ResourceDestructions)')
CHANGES=$(echo "$PLAN_JSON" | jq -r '.ResourceChanges')

if [ "$TOTAL" -gt 50 ]; then
  echo "WARNING: Large change set ($TOTAL resources)"
  echo "Additions: $(echo "$PLAN_JSON" | jq -r '.ResourceAdditions')"
  echo "Changes: $(echo "$PLAN_JSON" | jq -r '.ResourceChanges')"
  echo "Destructions: $(echo "$PLAN_JSON" | jq -r '.ResourceDestructions')"
fi
```

### 5. High-Risk Resource Checks

**Check for database changes:**

```bash
# Databases being destroyed
hcptf plan logs -id=run-abc123 | grep -i "db_instance\|rds\|database" | grep -i "destroy\|delete"

# Database configuration changes
hcptf plan logs -id=run-abc123 | grep -A 20 "aws_db_instance\|azurerm_sql_database\|google_sql_database_instance"
```

**Check for network changes:**

```bash
# VPC/Network changes
hcptf plan logs -id=run-abc123 | grep -i "vpc\|subnet\|route_table\|network" | grep -i "~\|+\|-"

# Security group changes
hcptf plan logs -id=run-abc123 | grep -A 10 "security_group"

# Load balancer changes
hcptf plan logs -id=run-abc123 | grep -i "load_balancer\|alb\|elb" | grep -i "~\|+\|-"
```

**Check for IAM/Permission changes:**

```bash
# IAM changes
hcptf plan logs -id=run-abc123 | grep -i "iam_role\|iam_policy\|role_assignment" | grep -i "~\|+\|-"

# Service principal changes
hcptf plan logs -id=run-abc123 | grep -i "service_account\|service_principal"
```

### 6. Validate Against Expectations

**Expected vs Actual Changes:**

```bash
# If you expected to add 1 S3 bucket, verify that's all that's changing
hcptf plan read -id=run-abc123 -output=json | jq '{
  additions: .ResourceAdditions,
  changes: .ResourceChanges,
  destructions: .ResourceDestructions
}'

# Then check the logs to see WHAT is being added
hcptf plan logs -id=run-abc123 | grep "will be created"
```

**Check specific resources:**

```bash
# Verify a specific resource is in the plan
hcptf plan logs -id=run-abc123 | grep "aws_s3_bucket.my_bucket"

# Check if any resources in production are affected
hcptf plan logs -id=run-abc123 | grep -i "prod\|production"

# Look for unexpected tag changes
hcptf plan logs -id=run-abc123 | grep -A 3 "tags"
```

### 7. Compare with Previous Plans

```bash
# Get last two runs
RUNS=$(hcptf run list -org=my-org -workspace=my-workspace -output=json | jq -r '.[0:2] | .[] | .ID')

# Compare plan stats
for run in $RUNS; do
  echo "Run: $run"
  hcptf plan read -id=$run -output=json | jq '{
    status: .Status,
    changes: {
      add: .ResourceAdditions,
      change: .ResourceChanges,
      destroy: .ResourceDestructions
    }
  }'
  echo ""
done
```

### 8. Safety Checklist

Before approving a plan, verify:

```bash
#!/bin/bash
# plan-safety-check.sh

RUN_ID="$1"

echo "=== PLAN SAFETY CHECKLIST ==="
echo "Run: $RUN_ID"
echo ""

# Get plan details
PLAN_JSON=$(hcptf plan read -id=$RUN_ID -output=json)

# 1. Check if plan has changes
HAS_CHANGES=$(echo "$PLAN_JSON" | jq -r '.HasChanges')
if [ "$HAS_CHANGES" = "false" ]; then
  echo "✓ No changes - safe to apply"
  exit 0
fi

# 2. Check destruction count
DESTRUCTIONS=$(echo "$PLAN_JSON" | jq -r '.ResourceDestructions')
echo "Destructions: $DESTRUCTIONS"
if [ "$DESTRUCTIONS" -gt 0 ]; then
  echo "⚠ WARNING: Resources will be destroyed"
  echo "Review these carefully:"
  hcptf plan logs -id=$RUN_ID | grep -B 2 "will be destroyed" | head -20
  echo ""
fi

# 3. Check for replacements
REPLACEMENTS=$(hcptf plan logs -id=$RUN_ID | grep -c "must be replaced" || echo "0")
echo "Replacements: $REPLACEMENTS"
if [ "$REPLACEMENTS" -gt 0 ]; then
  echo "⚠ WARNING: Resources will be replaced (destroy + create)"
  echo "This may cause downtime!"
  hcptf plan logs -id=$RUN_ID | grep -B 5 "must be replaced" | head -30
  echo ""
fi

# 4. Check for high-risk resources
echo "High-risk resource check:"
hcptf plan logs -id=$RUN_ID > /tmp/plan-$RUN_ID.txt

DB_CHANGES=$(grep -ic "db_instance\|database" /tmp/plan-$RUN_ID.txt || echo "0")
SG_CHANGES=$(grep -ic "security_group" /tmp/plan-$RUN_ID.txt || echo "0")
IAM_CHANGES=$(grep -ic "iam_role\|iam_policy" /tmp/plan-$RUN_ID.txt || echo "0")

echo "  Database changes: $DB_CHANGES"
echo "  Security group changes: $SG_CHANGES"
echo "  IAM changes: $IAM_CHANGES"

if [ "$DB_CHANGES" -gt 0 ] || [ "$SG_CHANGES" -gt 5 ] || [ "$IAM_CHANGES" -gt 3 ]; then
  echo "⚠ WARNING: High-risk changes detected"
fi

echo ""
echo "5. Total impact:"
echo "$PLAN_JSON" | jq '{
  additions: .ResourceAdditions,
  changes: .ResourceChanges,
  destructions: .ResourceDestructions,
  total: (.ResourceAdditions + .ResourceChanges + .ResourceDestructions)
}'

rm -f /tmp/plan-$RUN_ID.txt

echo ""
echo "Review complete. Carefully examine warnings before applying."
```

## Common Scenarios

### Scenario 1: Reviewing a Routine Update

```bash
# 1. Get plan summary
hcptf plan read -id=run-abc123

# 2. Quick check - should be mostly changes, few/no destructions
hcptf plan read -id=run-abc123 -output=json | jq '{
  additions: .ResourceAdditions,
  changes: .ResourceChanges,
  destructions: .ResourceDestructions
}'

# 3. Spot-check a few resources
hcptf plan logs -id=run-abc123 | grep "will be updated" | head -10

# 4. If everything looks good, apply
hcptf run apply -id=run-abc123 -comment="Reviewed and approved"
```

### Scenario 2: Investigating Unexpected Changes

```bash
# 1. Check what's changing
hcptf plan logs -id=run-abc123 > unexpected-plan.txt

# 2. Find resources you didn't expect to change
grep "~" unexpected-plan.txt | grep -i "prod"

# 3. Look at the specific changes
hcptf plan logs -id=run-abc123 | grep -A 30 "aws_instance.production_server"

# 4. Check configuration version to see what code changed
hcptf run read -id=run-abc123 -output=json | jq -r '.ConfigurationVersion.ID' | \
  xargs -I {} hcptf configversion read -id={}

# 5. Decision: if changes are unintended, discard and fix code
hcptf run discard -id=run-abc123 -comment="Unintended changes, fixing configuration"
```

### Scenario 3: Pre-Production Validation

```bash
# 1. Get plan
PLAN_JSON=$(hcptf plan read -id=run-abc123 -output=json)

# 2. Verify it matches expectations
EXPECTED_ADDS=5
ACTUAL_ADDS=$(echo "$PLAN_JSON" | jq -r '.ResourceAdditions')

if [ "$ACTUAL_ADDS" -ne "$EXPECTED_ADDS" ]; then
  echo "ERROR: Expected $EXPECTED_ADDS additions, got $ACTUAL_ADDS"
  hcptf plan logs -id=run-abc123 | grep "will be created"
  exit 1
fi

# 3. Check no unexpected deletions
DESTRUCTIONS=$(echo "$PLAN_JSON" | jq -r '.ResourceDestructions')
if [ "$DESTRUCTIONS" -gt 0 ]; then
  echo "ERROR: Unexpected destructions"
  hcptf plan logs -id=run-abc123 | grep "will be destroyed"
  exit 1
fi

# 4. Validate specific resources
hcptf plan logs -id=run-abc123 | grep "aws_s3_bucket.data" || {
  echo "ERROR: Expected S3 bucket not in plan"
  exit 1
}

echo "✓ Plan validation passed"
```

### Scenario 4: Database Change Review

```bash
# 1. Check plan for database changes
hcptf plan logs -id=run-abc123 | grep -i "db_instance" > db-changes.txt

# 2. Review what's changing
cat db-changes.txt

# 3. Check if it's a replacement (BAD - causes downtime)
if grep -q "must be replaced" db-changes.txt; then
  echo "⚠ CRITICAL: Database will be replaced!"
  echo "This will:"
  echo "  - Delete the existing database"
  echo "  - Create a new database"
  echo "  - Cause data loss and downtime"
  echo ""
  echo "Reason for replacement:"
  hcptf plan logs -id=run-abc123 | grep -B 10 -A 5 "aws_db_instance.*must be replaced"

  # Consider using blue/green deployment instead
  exit 1
fi

# 4. Check specific attributes changing
echo "Database configuration changes:"
hcptf plan logs -id=run-abc123 | grep -A 30 "aws_db_instance" | \
  grep -E "~|instance_class|engine_version|storage"

# 5. If it's just a version upgrade or scaling, might be OK
echo ""
echo "Review complete - check if changes are acceptable"
```

## Tips and Best Practices

1. **Always review before applying**: Never blindly apply without checking the plan
2. **Watch for replacements**: These cause downtime for stateful resources
3. **Check destructions carefully**: Especially for databases, stateful resources
4. **Validate scope matches intent**: Small config change shouldn't affect 50 resources
5. **Use safety checks in automation**: Script the common validation patterns
6. **Compare with previous runs**: Look for anomalies in change patterns
7. **Check production tags**: Verify production resources aren't unexpectedly changing
8. **Review security changes**: Security groups, IAM roles, network config
9. **Save plan output**: Keep records of what was planned vs. applied
10. **Understand replacement reasons**: Check why Terraform requires replace

## Troubleshooting

**Plan shows no changes but you expect changes:**
- Check if configuration was actually modified
- Verify VCS connection triggered a run with new code
- Review the configuration version: `hcptf run read -id=run-xyz -output=json | jq '.ConfigurationVersion'`

**Plan shows unexpected destructions:**
- Check if resources were removed from configuration
- Verify state hasn't drifted (use drift detection skill)
- Review recent commits to configuration

**Plan is too large to review:**
- Focus on destructions first: `hcptf plan logs -id=run-xyz | grep "will be destroyed"`
- Then replacements: `grep "must be replaced"`
- Filter by resource type: `grep "aws_db_instance"`

**Can't find recent plan:**
- List runs: `hcptf run list -org=my-org -workspace=my-workspace`
- Check run status: plans are only available for runs in certain states
- Verify workspace name: `hcptf workspace list -org=my-org`

## Related Commands

- `hcptf run list` - List runs for a workspace
- `hcptf run read` - Get run details including plan ID
- `hcptf plan read` - Get plan summary and statistics
- `hcptf plan logs` - View full plan text output
- `hcptf run apply` - Apply a plan after review
- `hcptf run discard` - Discard a plan if issues found
- `hcptf configversion read` - See what code version is being applied

## Integration with Other Skills

- **drift**: Use plan-analyzer after fixing drift to verify remediation
- **state-analyzer**: Compare planned changes with current state
- **policy-compliance**: Plans must pass policy checks before review
- **version-upgrades**: Validate upgrade plans don't cause unexpected changes

## Agent Considerations

When building agents that analyze plans:

1. **Prioritize Safety**: Flag destructions and replacements prominently
2. **Explain Impact**: Help users understand what each change means
3. **Highlight Risks**: Call out high-risk resources (databases, networks, IAM)
4. **Validate Intent**: Check if changes match user's stated goals
5. **Provide Context**: Show why resources are being replaced
6. **Suggest Actions**: "Review database changes before applying"
7. **Compare Patterns**: "This is more changes than usual for this workspace"
8. **Check Dependencies**: Understand cascading effects
9. **Link to Code**: Connect plan changes to configuration changes
10. **Document Decision**: Record why plan was approved/rejected
