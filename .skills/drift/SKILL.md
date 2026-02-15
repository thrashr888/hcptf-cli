# Drift Detection and Resolution Skill

## Overview
This skill helps investigate and resolve infrastructure drift in HCP Terraform workspaces. Drift occurs when actual infrastructure state diverges from the Terraform configuration.

## Prerequisites
- HCP Terraform Plus or Enterprise (health assessments feature)
- Health assessments enabled in workspace settings
- Authenticated with `hcptf` CLI (`TFE_TOKEN` or `~/.terraform.d/credentials.terraform.io`)

## Core Concepts

**Drift Types:**
- **delete**: Resource was deleted outside Terraform - usually needs `terraform apply` to recreate
- **update**: Resource was modified outside Terraform - may need code update or apply
- **create**: Resource was created outside Terraform - may need import or removal

**Resolution Strategies:**
1. **Apply fixes infrastructure**: Run `terraform apply` to make infrastructure match code
2. **Update code**: Modify Terraform configuration to match actual infrastructure state
3. **Update variables**: Change variable values if drift is due to outdated variable values
4. **Import resources**: Import manually created resources into Terraform state

## Workflow

### 1. Find Drifted Workspaces

Use Explorer API to find workspaces with drift:

```bash
# Find all drifted workspaces in organization
hcptf explorer query -org=<org-name> -type=workspaces \
  -fields=workspace-name,drifted,resources-drifted,resources-undrifted \
  | grep -v "false"

# Sort by most drifted resources
hcptf explorer query -org=<org-name> -type=workspaces \
  -sort=-resources-drifted -page-size=20
```

### 2. View Drift Details

Get detailed drift information for a workspace:

```bash
# View drift and check results
hcptf assessmentresult list -org=<org-name> -name=<workspace-name>
```

This shows:
- Which resources drifted
- What actions are needed (delete, update, create)
- Which attributes changed (before/after values)
- Terraform check failures (continuous validation)

### 3. Understand the Configuration

**Get workspace VCS information:**

```bash
# View workspace details including VCS repo
hcptf explorer query -org=<org-name> -type=workspaces \
  -fields=workspace-name,vcs-repo-identifier \
  -filter="workspace-name:<workspace>"
```

**Get commit information from the run:**

```bash
# Get current run ID
RUN_ID=$(hcptf workspace read -org=<org> -name=<workspace> -output=json | jq -r '.CurrentRunID')

# Get configuration version
CONFIG_VERSION=$(hcptf run show -id=$RUN_ID -output=json | jq -r '.ConfigurationVersionID')

# View configuration version details (includes VCS commit info if available)
hcptf configversion read -id=$CONFIG_VERSION
```

**For VCS-backed workspaces**, this provides:
- Branch name
- Commit SHA
- Commit URL (direct link to GitHub/GitLab commit)
- Repo identifier (e.g., `username/repo`)

### 4. Analyze and Decide

Review the drift details to determine the appropriate action:

**Decision Matrix:**

| Action Type | Common Cause | Typical Resolution |
|-------------|-------------|-------------------|
| `delete` | Manual deletion, autoscaling | Apply to recreate |
| `update` with infrastructure values | Manual changes, console edits | Apply to revert OR update code to match |
| `update` with computed values | Normal operation (IPs, timestamps) | Update code/variables to match |
| Tags added | Tagging automation | Update code to include tags |
| Check failures | Cert expiration, thresholds | Varies by check type |

**Questions to ask:**
- Is the current infrastructure state correct?
- Should the code be updated to match reality?
- Is this drift expected (computed values, auto-scaling)?
- Do variables need updating?

### 5. Take Action

**Option A: Fix infrastructure (apply)**

```bash
# Create a run to fix drift
hcptf run create -org=<org> -name=<workspace> \
  -message="Fix drift: [describe what's being fixed]"

# Monitor the run
hcptf run show -id=<run-id>
```

**Option B: Update code**

1. Clone the repository (use VCS info from step 3)
2. Checkout the branch that's connected to the workspace
3. Update Terraform files to match actual infrastructure
4. Commit and push - this triggers a new run automatically

**Option C: Update variables**

```bash
# List current variables
hcptf variable list -org=<org> -workspace=<workspace>

# Update a variable
hcptf variable update -org=<org> -workspace=<workspace> \
  -key=<variable-name> -value=<new-value>
```

**Option D: Import resources (for manually created resources)**

```bash
# This requires terraform CLI access to the workspace
terraform import <resource_address> <resource_id>
```

### 6. Verify Resolution

After taking action, verify drift is resolved:

```bash
# Check if drift is cleared (may take a few minutes for next health check)
hcptf explorer query -org=<org> -type=workspaces \
  -filter="workspace-name:<workspace>" \
  -fields=workspace-name,drifted,resources-drifted

# View latest assessment result
hcptf assessmentresult list -org=<org> -name=<workspace>
```

## Common Scenarios

### Scenario 1: EC2 Instance Deleted (Autoscaling)

```bash
# 1. Identify drift
hcptf assessmentresult list -org=my-org -name=my-workspace
# Shows: aws_instance.app - Action: delete

# 2. Decision: Recreate the instance
# 3. Apply fix
hcptf run create -org=my-org -name=my-workspace \
  -message="Recreate deleted EC2 instance"
```

### Scenario 2: Route53 Record IP Changed

```bash
# 1. View drift
hcptf assessmentresult list -org=my-org -name=my-workspace
# Shows: aws_route53_record.app - Action: update
#        records: ["1.2.3.4"] -> ["5.6.7.8"]

# 2. Decision: Is new IP correct?
#    - If YES: Update variable or code
#    - If NO: Run apply to revert

# 3a. If updating variable:
hcptf variable update -org=my-org -workspace=my-workspace \
  -key=app_ip -value="5.6.7.8"

# 3b. If reverting infrastructure:
hcptf run create -org=my-org -name=my-workspace \
  -message="Revert Route53 record to configured IP"
```

### Scenario 3: Tags Added by Automation

```bash
# 1. View drift
hcptf assessmentresult list -org=my-org -name=my-workspace
# Shows: multiple resources - Action: update
#        tags: {} -> {"Environment": "prod", "Owner": "team-a"}

# 2. Decision: Tags are correct, update code to include them
# 3. Update Terraform code to include tags in resource definitions
# 4. Commit and push (triggers auto-run if VCS-backed)
```

### Scenario 4: Certificate Expiring (Check Failure)

```bash
# 1. View checks
hcptf assessmentresult list -org=my-org -name=my-workspace
# Shows: TERRAFORM CHECK RESULTS
#        tls_self_signed_cert.app - FAIL
#        • Certificate will expire in less than 4 hours

# 2. Decision: Need to regenerate certificate
# 3. Update code to generate new cert or increase validity period
# 4. Apply changes
hcptf run create -org=my-org -name=my-workspace \
  -message="Regenerate expiring certificate"
```

## Tips and Best Practices

1. **Regular monitoring**: Check drift regularly using Explorer queries
2. **Investigate before fixing**: Always understand WHY drift occurred
3. **Document decisions**: Include clear messages when creating runs
4. **Use continuous validation**: Enable Terraform checks to catch issues early
5. **Automate remediation**: For known drift patterns, consider automated fixes
6. **Review permissions**: Drift often indicates someone has console access they shouldn't
7. **Update documentation**: If updating code to match infrastructure, document why

## Troubleshooting

**No assessment results found:**
- Health assessments may not be enabled in workspace settings
- Feature requires HCP Terraform Plus or Enterprise
- Enable via UI or `hcptf workspace update`

**401 Unauthorized on assessment results:**
- Check your authentication token is valid
- Ensure you have at least read access to the workspace

**Drift shows but no resources listed:**
- May be computed attribute changes (not shown in detail)
- Check the log output URL for full details

**Apply doesn't fix drift:**
- Drift may be due to external automation re-applying changes
- Consider updating code to match the pattern
- Investigate what's causing the external changes

## Related Commands

- `hcptf explorer query` - Query workspaces and resources
- `hcptf assessmentresult list` - View drift and check details
- `hcptf workspace read` - Get workspace details
- `hcptf run create` - Create a run to fix drift
- `hcptf run show` - Monitor run progress
- `hcptf variable list/update` - Manage workspace variables
- `hcptf configversion read` - Get VCS commit information

## Agent Considerations

When building agents that handle drift:

1. **Assess severity**: Check if drift is critical or expected
2. **Get approval**: Don't auto-fix drift without user review
3. **Provide context**: Explain what changed and why
4. **Link to code**: Use VCS URLs to show relevant code
5. **Suggest actions**: Present options with pros/cons
6. **Verify resolution**: Confirm drift is cleared after action
