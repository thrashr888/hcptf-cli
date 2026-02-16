---
name: state-analyzer
description: Analyze Terraform state files and provide improvement recommendations. Use when investigating state health, optimizing infrastructure, reviewing security posture, identifying cost savings, or auditing resource configurations in HCP Terraform workspaces.
---

# State Analyzer Skill

## Overview

This skill analyzes Terraform state files to identify security issues, performance bottlenecks, cost optimization opportunities, and best practice violations. It provides actionable recommendations for improving your infrastructure-as-code.

## Prerequisites

- HCP Terraform account with workspace access
- Authenticated with `hcptf` CLI (`TFE_TOKEN` or `~/.terraform.d/credentials.terraform.io`)
- Read access to workspace state files
- For local analysis: Terraform state file (`terraform.tfstate`)

## Core Concepts

**State File Analysis Areas:**

1. **Security**: Exposed secrets, insecure configurations, open security groups
2. **Cost Optimization**: Over-provisioned resources, unused resources, cheaper alternatives
3. **Performance**: Resource sizing, network topology, caching opportunities
4. **Best Practices**: Naming conventions, tagging, modularity, dependencies
5. **Compliance**: Policy violations, required tags, encryption standards
6. **State Health**: Large state files, orphaned resources, circular dependencies

**Common Issues Detected:**

- Hardcoded sensitive values in resource attributes
- Resources with public access (S3 buckets, security groups, databases)
- Over-provisioned compute instances (large instance types with low utilization)
- Missing or inconsistent resource tags
- Deprecated resource types or configurations
- Unused resources (elastic IPs, volumes, load balancers)
- Inefficient network architectures
- Non-compliance with organizational standards

## Workflow

### 1. Get State File

**Option A: Download from HCP Terraform Workspace**

```bash
# Download current state to file
hcptf state download -org=my-org -workspace=my-workspace -output=state.json

# Or print to stdout for piping
hcptf state download -org=my-org -workspace=my-workspace | jq '.'

# Direct analysis without saving to file
hcptf state download -org=my-org -workspace=my-workspace | jq '.resources | length'
```

**Option B: Use Local State File**

```bash
# If working with local Terraform
terraform show -json > state.json
```

### 2. Parse State Structure

Examine the state file structure to understand what you're working with:

```bash
# Get high-level statistics (from file)
jq '{
  terraform_version: .terraform_version,
  serial: .serial,
  lineage: .lineage,
  resource_count: (.resources | length),
  output_count: (.outputs | length)
}' state.json

# Or directly from workspace (no file needed)
hcptf state download -org=my-org -workspace=my-workspace | jq '{
  terraform_version,
  serial,
  resource_count: (.resources | length),
  output_count: (.outputs | length)
}'

# List all resources by type
jq -r '.resources[] | "\(.type) - \(.name)"' state.json | sort | uniq -c

# Or directly from workspace
hcptf state download -org=my-org -workspace=my-workspace | \
  jq -r '.resources[] | "\(.type) - \(.name)"' | sort | uniq -c
```

### 3. Security Analysis

**Check for exposed secrets:**

```bash
# Search for common secret patterns in resource attributes
jq -r '.resources[] |
  select(.instances[0].attributes |
    (. | tostring) | test("password|secret|key|token"; "i")
  ) |
  "\(.type).\(.name)"' state.json
```

**Identify publicly accessible resources:**

```bash
# Find S3 buckets with public access
jq -r '.resources[] |
  select(.type == "aws_s3_bucket") |
  select(.instances[0].attributes.acl == "public-read" or
         .instances[0].attributes.acl == "public-read-write") |
  "\(.name) - Public ACL: \(.instances[0].attributes.acl)"' state.json

# Find security groups with 0.0.0.0/0 access
jq -r '.resources[] |
  select(.type == "aws_security_group") |
  select(.instances[0].attributes.ingress[]? | .cidr_blocks[]? == "0.0.0.0/0") |
  "\(.name) - Open to internet"' state.json

# Find RDS instances without encryption
jq -r '.resources[] |
  select(.type == "aws_db_instance") |
  select(.instances[0].attributes.storage_encrypted == false) |
  "\(.name) - Unencrypted storage"' state.json
```

### 4. Cost Optimization Analysis

**Identify over-provisioned resources:**

```bash
# List large EC2 instances
jq -r '.resources[] |
  select(.type == "aws_instance") |
  select(.instances[0].attributes.instance_type |
    test("xlarge|2xlarge|4xlarge|8xlarge")) |
  "\(.name) - \(.instances[0].attributes.instance_type)"' state.json

# Find provisioned IOPS volumes
jq -r '.resources[] |
  select(.type == "aws_ebs_volume") |
  select(.instances[0].attributes.iops != null) |
  "\(.name) - IOPS: \(.instances[0].attributes.iops)"' state.json
```

**Detect unused resources:**

```bash
# Find unattached EBS volumes
jq -r '.resources[] |
  select(.type == "aws_ebs_volume") |
  select(.instances[0].attributes.attachments | length == 0) |
  "\(.name) - Unattached volume"' state.json

# Find elastic IPs not associated with instances
jq -r '.resources[] |
  select(.type == "aws_eip") |
  select(.instances[0].attributes.instance == null or
         .instances[0].attributes.instance == "") |
  "\(.name) - Unassociated EIP"' state.json
```

### 5. Best Practices Analysis

**Check tagging consistency:**

```bash
# Find resources without required tags
REQUIRED_TAGS=("Environment" "Owner" "Project")

jq -r --argjson tags '["Environment","Owner","Project"]' '
  .resources[] |
  select(.instances[0].attributes.tags) |
  select([.instances[0].attributes.tags | keys[] as $k |
    select($tags | index($k) | not)] | length > 0) |
  "\(.type).\(.name) - Missing tags: \(
    $tags - (.instances[0].attributes.tags | keys) | join(", ")
  )"' state.json
```

**Check naming conventions:**

```bash
# Identify resources not following naming patterns
# Example: expecting format like "{env}-{project}-{resource}-{name}"
jq -r '.resources[] |
  select(.name | test("^[a-z]+(-[a-z0-9]+)*$") | not) |
  "\(.type).\(.name) - Non-standard naming"' state.json
```

**Find deprecated resource types:**

```bash
# AWS example - check for deprecated resources
jq -r '.resources[] |
  select(.type | test(
    "aws_elasticache_security_group|aws_db_security_group|aws_redshift_security_group"
  )) |
  "\(.type).\(.name) - Deprecated resource type"' state.json
```

### 6. Performance Analysis

**Analyze resource dependencies:**

```bash
# Count dependencies per resource
jq -r '.resources[] |
  "\(.type).\(.name) - Dependencies: \(
    if .instances[0].dependencies then
      (.instances[0].dependencies | length)
    else 0 end
  )"' state.json | sort -t: -k2 -nr | head -20
```

**Identify potential bottlenecks:**

```bash
# Find resources with many dependencies (potential bottleneck)
jq -r '.resources[] |
  select(.instances[0].dependencies | length > 5) |
  "\(.type).\(.name) - \(.instances[0].dependencies | length) dependencies"' state.json
```

### 7. State Health Analysis

**Check state file size and complexity:**

```bash
# State file statistics
jq '{
  total_resources: (.resources | length),
  total_outputs: (.outputs | length),
  resource_types: (.resources | group_by(.type) |
    map({type: .[0].type, count: length}) |
    sort_by(.count) | reverse | .[0:10]
  ),
  largest_resources: (.resources |
    map({
      name: "\(.type).\(.name)",
      size: (.instances[0].attributes | tostring | length)
    }) |
    sort_by(.size) | reverse | .[0:5]
  )
}' state.json
```

**Detect potential issues:**

```bash
# Find resources with many instances (potential split needed)
jq -r '.resources[] |
  select(.instances | length > 1) |
  "\(.type).\(.name) - \(.instances | length) instances"' state.json

# Check for very large individual resources
jq -r '.resources[] |
  select((.instances[0].attributes | tostring | length) > 50000) |
  "\(.type).\(.name) - Large resource (consider splitting)"' state.json
```

### 8. Generate Recommendations Report

Create a comprehensive analysis report:

```bash
# Save this as analyze_state.sh
#!/bin/bash
STATE_FILE="$1"

echo "# Terraform State Analysis Report"
echo "Generated: $(date)"
echo ""

echo "## Summary"
jq '{
  terraform_version: .terraform_version,
  total_resources: (.resources | length),
  total_outputs: (.outputs | length),
  providers: (.resources | map(.provider) | unique)
}' "$STATE_FILE"

echo ""
echo "## Security Issues"
echo "### Public S3 Buckets"
jq -r '.resources[] |
  select(.type == "aws_s3_bucket") |
  select(.instances[0].attributes.acl | test("public")) |
  "- \(.name): \(.instances[0].attributes.acl)"' "$STATE_FILE"

echo ""
echo "### Open Security Groups"
jq -r '.resources[] |
  select(.type == "aws_security_group") |
  select(.instances[0].attributes.ingress[]? | .cidr_blocks[]? == "0.0.0.0/0") |
  "- \(.name): Allows access from 0.0.0.0/0"' "$STATE_FILE"

echo ""
echo "## Cost Optimization"
echo "### Large Instances"
jq -r '.resources[] |
  select(.type == "aws_instance") |
  select(.instances[0].attributes.instance_type | test("xlarge")) |
  "- \(.name): \(.instances[0].attributes.instance_type)"' "$STATE_FILE"

echo ""
echo "### Unattached Resources"
jq -r '.resources[] |
  select(.type == "aws_ebs_volume") |
  select(.instances[0].attributes.attachments | length == 0) |
  "- \(.name): Unattached EBS volume"' "$STATE_FILE"

echo ""
echo "## Best Practices"
echo "### Missing Tags"
jq -r '.resources[] |
  select(.instances[0].attributes.tags?) |
  select(.instances[0].attributes.tags.Environment? | not) |
  "- \(.type).\(.name): Missing Environment tag"' "$STATE_FILE"
```

Usage:

```bash
chmod +x analyze_state.sh
./analyze_state.sh state.json > analysis_report.md
```

## Common Analysis Scenarios

### Scenario 1: Security Audit Before Production

```bash
# 1. Download production workspace state
hcptf state download -org=my-org -workspace=prod-app -output=prod-state.json

# Or analyze directly without saving
# 2. Run comprehensive security checks
echo "=== Public Access Check ==="
jq -r '.resources[] |
  select(.type | test("aws_s3_bucket|aws_security_group|aws_db_instance")) |
  select(
    (.instances[0].attributes.acl? | test("public")) or
    (.instances[0].attributes.publicly_accessible? == true) or
    (.instances[0].attributes.ingress[]?.cidr_blocks[]? == "0.0.0.0/0")
  ) |
  "\(.type).\(.name) - Public access detected"' prod-state.json

echo ""
echo "=== Encryption Check ==="
jq -r '.resources[] |
  select(.type | test("aws_db_instance|aws_ebs_volume|aws_s3_bucket")) |
  select(
    (.instances[0].attributes.storage_encrypted? == false) or
    (.instances[0].attributes.encrypted? == false) or
    (.instances[0].attributes.server_side_encryption_configuration? | length == 0)
  ) |
  "\(.type).\(.name) - Encryption disabled"' prod-state.json

# 3. Document findings and remediate
```

### Scenario 2: Cost Optimization Review

```bash
# 1. Get state from all environments
for env in dev staging prod; do
  hcptf my-org ${env}-app statefile download -output=${env}-state.json
done

# 2. Compare resource sizes across environments
for env in dev staging prod; do
  echo "=== $env Environment ==="
  jq -r '.resources[] |
    select(.type == "aws_instance") |
    "\(.name): \(.instances[0].attributes.instance_type)"' ${env}-state.json
done

# 3. Identify over-provisioned dev/staging
# Dev should typically use smaller instance types than prod

# 4. Generate cost savings recommendations
```

### Scenario 3: Compliance Check

```bash
# 1. Download workspace state
hcptf my-org compliance-workspace statefile download -output=state.json

# 2. Check required tags
REQUIRED_TAGS='["Environment","Owner","CostCenter","Project"]'

jq --argjson tags "$REQUIRED_TAGS" -r '
  .resources[] |
  select(.instances[0].attributes.tags?) |
  . as $resource |
  $tags - (.instances[0].attributes.tags | keys) |
  select(length > 0) |
  "\($resource.type).\($resource.name) - Missing: \(. | join(", "))"
' state.json > missing-tags.txt

# 3. Check encryption requirements
jq -r '.resources[] |
  select(.type | test("aws_s3_bucket|aws_ebs_volume|aws_db_instance")) |
  select(
    (.instances[0].attributes.encrypted? == false) or
    (.instances[0].attributes.storage_encrypted? == false) or
    (.instances[0].attributes.server_side_encryption_configuration? | length == 0)
  ) |
  "\(.type).\(.name)"' state.json > unencrypted-resources.txt

# 4. Generate compliance report
```

### Scenario 4: Pre-Migration State Assessment

```bash
# 1. Analyze current state before migration
hcptf my-org legacy-workspace statefile download -output=current-state.json

# 2. Identify resources that need special handling
echo "=== Stateful Resources (need backup/migration plan) ==="
jq -r '.resources[] |
  select(.type | test("aws_db_instance|aws_dynamodb_table|aws_s3_bucket")) |
  "\(.type).\(.name)"' current-state.json

# 3. Check for deprecated resources
echo "=== Deprecated Resources (need updating) ==="
jq -r '.resources[] |
  select(.type | test("aws_db_security_group|aws_elasticache_security_group")) |
  "\(.type).\(.name) - Deprecated, use VPC security groups"' current-state.json

# 4. Identify complex dependencies
echo "=== Highly Dependent Resources (migration complexity) ==="
jq -r '.resources[] |
  select(.instances[0].dependencies? | length > 5) |
  "\(.type).\(.name) - \(.instances[0].dependencies | length) dependencies"' current-state.json

# 5. Create migration plan based on findings
```

## Tips and Best Practices

1. **Regular Analysis**: Run state analysis monthly or before major changes
2. **Baseline First**: Establish a baseline analysis to track improvements
3. **Automate Checks**: Integrate state analysis into CI/CD pipelines
4. **Prioritize Findings**: Focus on security issues first, then cost, then best practices
5. **Compare Environments**: Analyze differences between dev/staging/prod for consistency
6. **Track Metrics**: Monitor state file size and resource count over time
7. **Version History**: Compare state across versions to identify trends
8. **Document Exceptions**: Some findings may be intentional - document these
9. **Team Review**: Share analysis results with team for collective improvement
10. **Action Items**: Convert findings into actionable Terraform changes

## Advanced Analysis Patterns

### Multi-Workspace Comparison

```bash
# Compare resource counts across workspaces
for ws in $(hcptf explorer query -org=my-org -type=workspaces \
            -fields=workspace-name -output=json | jq -r '.[].workspace_name'); do
  count=$(hcptf my-org $ws statefile -output=json | jq '.Resources | length')
  echo "$ws: $count resources"
done | sort -t: -k2 -nr
```

### Trend Analysis

```bash
# Get historical state versions and track growth
hcptf my-org my-workspace statefile list | head -10 | while read -r version; do
  hcptf my-org my-workspace statefile download \
    -state-version=$version -output=state-${version}.json

  count=$(jq '.resources | length' state-${version}.json)
  date=$(jq -r '.created_at' state-${version}.json)
  echo "$date: $count resources"
done
```

### Cross-Provider Analysis

```bash
# Identify which providers are used
jq -r '.resources | group_by(.provider) |
  map({provider: .[0].provider, count: length}) |
  sort_by(.count) | reverse |
  .[] | "\(.provider): \(.count) resources"' state.json
```

## Troubleshooting

**Cannot download state file:**
- Verify workspace permissions (need at least read access)
- Check authentication: `hcptf organizations list`
- Workspace may not have any state yet (no successful applies)

**jq command not found:**
```bash
# Install jq for JSON parsing
brew install jq  # macOS
apt-get install jq  # Ubuntu/Debian
```

**State file too large to analyze:**
- Use jq streaming parser: `jq -c '.resources[]' state.json | while read -r res; do ...`
- Focus on specific resource types rather than full analysis
- Consider splitting large workspaces

**False positives in security scan:**
- Review each finding - some may be intentional (bastion hosts, public websites)
- Maintain an exceptions list for known-good configurations
- Update queries to exclude approved patterns

## Related Commands

- `hcptf statefile download` - Download workspace state
- `hcptf statefile list` - List state versions
- `hcptf runs list` - View recent runs that modified state
- `hcptf workspace read` - Get workspace configuration
- `terraform show -json` - Convert binary state to JSON

## Integration with Other Skills

- **drift**: Use state-analyzer to understand baseline before investigating drift
- **policy-compliance**: Combine with Sentinel policies for enforcement
- **version-upgrades**: Identify deprecated resources before provider upgrades
- **workspace-to-stack**: Assess state complexity before converting to stacks

## Agent Considerations

When building agents that analyze state:

1. **Permissions**: Always verify read access before downloading state
2. **Sensitive Data**: State files contain secrets - handle securely, never log
3. **Context Matters**: Some findings require business context to interpret
4. **Prioritization**: Present high-severity issues (security) before low (naming)
5. **Actionable Output**: Provide specific remediation steps, not just findings
6. **Performance**: Large states may need streaming or sampling
7. **Versioning**: Specify which state version was analyzed
8. **Comparison**: Offer to compare with previous analysis for trends
9. **Export Options**: Support multiple output formats (JSON, Markdown, HTML)
10. **Follow-up**: Suggest related skills or next steps based on findings
