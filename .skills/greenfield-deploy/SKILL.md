---
name: greenfield-deploy
description: Set up a brand-new infrastructure project from scratch using HCP Terraform. Use when starting a new project from zero, creating workspaces with VCS integration, configuring variables and credentials, running first deployments, or onboarding a new team to HCP Terraform.
---

# Greenfield Project Deployment Skill

## Overview

This skill guides setting up a brand-new infrastructure project from scratch using HCP Terraform. It covers the full end-to-end workflow: creating the organization structure, writing Terraform configuration, setting up a workspace with VCS integration, configuring variables and credentials, running the first deployment, and verifying everything is working. This is the "day zero" workflow for teams starting fresh with HCP Terraform.

## Prerequisites

- Authenticated with `hcptf` CLI (`TFE_TOKEN` or `~/.terraform.d/credentials.terraform.io`)
- An HCP Terraform organization (or permissions to create one)
- A VCS account (GitHub, GitLab, Bitbucket, or Azure DevOps) with an OAuth connection configured
- Cloud provider credentials (AWS, Azure, GCP, etc.) for the target infrastructure

## Core Concepts

**HCP Terraform Resource Hierarchy:**

- **Organization**: Top-level container for all resources
- **Project**: Groups related workspaces (e.g., "platform", "application")
- **Workspace**: Manages a single Terraform configuration and its state
- **Run**: A plan/apply cycle within a workspace

**Deployment Approaches:**

1. **VCS-driven** (recommended): Connect workspace to a Git repo; pushes trigger runs automatically
2. **CLI-driven**: Use `terraform` CLI locally to trigger runs in HCP Terraform
3. **API-driven**: Use the API or `hcptf` CLI to upload config and trigger runs

**Variable Types:**

- **Terraform variables**: Passed to your `.tf` configuration (key/value pairs)
- **Environment variables**: Set in the run environment (e.g., `AWS_ACCESS_KEY_ID`)

## Workflow

### 1. Plan the Project Structure

Before creating anything, decide on your organization:

```
Organization: my-company
  Project: platform-infrastructure
    Workspace: networking-prod     (VPC, subnets, DNS)
    Workspace: networking-staging
  Project: application-services
    Workspace: api-service-prod    (ECS, ALB, RDS)
    Workspace: api-service-staging
```

### 2. Verify Authentication and Organization

```bash
# Login to HCP Terraform
hcptf login

# Verify authentication
hcptf account show

# List your organizations
hcptf organization list

# View organization details
hcptf <org>
```

### 3. Create Project

```bash
# List existing projects
hcptf project list -org=<org>

# Create a new project
hcptf project create -org=<org> \
  -name=<project-name> \
  -description="Infrastructure for <purpose>"

# Verify
hcptf project list -org=<org>
```

### 4. Write Terraform Configuration

Create your Terraform configuration in a Git repository. Minimum files:

**versions.tf:**

```hcl
terraform {
  required_version = ">= 1.10.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  cloud {
    organization = "my-org"
    workspaces {
      name = "my-workspace"
    }
  }
}
```

**main.tf:**

```hcl
provider "aws" {
  region = var.region
}

# Your resources here
resource "aws_s3_bucket" "app" {
  bucket = "${var.project_name}-${var.environment}"

  tags = {
    Environment = var.environment
    ManagedBy   = "Terraform"
    Project     = var.project_name
  }
}
```

**variables.tf:**

```hcl
variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
}
```

**outputs.tf:**

```hcl
output "bucket_arn" {
  description = "ARN of the S3 bucket"
  value       = aws_s3_bucket.app.arn
}
```

Push this to your VCS provider:

```bash
git init
git add .
git commit -m "Initial Terraform configuration"
git remote add origin https://github.com/<org>/<repo>.git
git push -u origin main
```

### 5. Create the Workspace

**Find your OAuth token for VCS:**

```bash
# List OAuth clients and tokens
hcptf oauthclient list -org=<org>
hcptf oauthtoken list -org=<org>
# Note the OAuth token ID (ot-xxx)
```

**Create the workspace:**

```bash
# Create workspace with VCS connection
hcptf workspace create -org=<org> -name=<workspace-name> \
  -auto-apply=false \
  -terraform-version=1.10.0 \
  -description="<purpose> - <environment>"

# Verify workspace was created
hcptf <org> <workspace>
```

**For VCS-connected workspaces** (recommended), configure VCS via the UI or API since the workspace create command has limited VCS flags. Alternatively, set up the `cloud` block in your Terraform config and use `terraform init` locally.

### 6. Configure Variables

**Set Terraform variables:**

```bash
# Required variables
hcptf variable create -org=<org> -workspace=<workspace> \
  -key=environment -value=production

hcptf variable create -org=<org> -workspace=<workspace> \
  -key=project_name -value=my-app

hcptf variable create -org=<org> -workspace=<workspace> \
  -key=region -value=us-east-1
```

**Set cloud provider credentials (environment variables):**

```bash
# AWS credentials
hcptf variable create -org=<org> -workspace=<workspace> \
  -key=AWS_ACCESS_KEY_ID -value=<access-key> -category=env -sensitive

hcptf variable create -org=<org> -workspace=<workspace> \
  -key=AWS_SECRET_ACCESS_KEY -value=<secret-key> -category=env -sensitive

hcptf variable create -org=<org> -workspace=<workspace> \
  -key=AWS_DEFAULT_REGION -value=us-east-1 -category=env
```

**For OIDC-based authentication (recommended for production):**

```bash
# Set up AWS OIDC trust
hcptf awsoidc create -org=<org> -name=aws-prod \
  -role-arn=arn:aws:iam::123456789:role/terraform-cloud
```

**Verify variables:**

```bash
hcptf variable list -org=<org> -workspace=<workspace>
```

### 7. Trigger the First Run

```bash
# Create the first run
hcptf workspace run create -org=<org> -name=<workspace> \
  -message="Initial infrastructure deployment"

# Get the run ID
hcptf <org> <workspace> runs

# View run details and plan output
RUN_ID=$(hcptf <org> <workspace> -output=json | jq -r '.CurrentRunID')
hcptf <org> <workspace> runs $RUN_ID show
```

### 8. Review and Apply

```bash
# Review the plan
hcptf <org> <workspace> runs $RUN_ID show
# Look for: Status, ResourceAdditions, ResourceChanges, ResourceDestructions

# If the plan looks good, apply
hcptf <org> <workspace> runs $RUN_ID apply -comment="Initial deployment approved"

# Monitor apply progress
hcptf <org> <workspace> runs $RUN_ID show
# Wait for Status: applied
```

### 9. Verify Deployment

```bash
# Check workspace state
hcptf state outputs -org=<org> -workspace=<workspace>

# View state versions
hcptf state list -org=<org> -workspace=<workspace>

# Verify workspace status
hcptf <org> <workspace>
# Look for: CurrentRunStatus: applied
```

### 10. Set Up Notifications (Optional)

```bash
# Slack notifications
hcptf notification create -org=<org> -workspace=<workspace> \
  -name=slack-alerts \
  -destination-type=slack \
  -url=https://hooks.slack.com/services/xxx/yyy/zzz \
  -triggers=run:completed,run:errored

# Verify
hcptf notification list -org=<org> -workspace=<workspace>
```

### 11. Configure Team Access (Optional)

```bash
# List teams
hcptf team list -org=<org>

# Grant team access to workspace
hcptf teamaccess create -org=<org> -workspace=<workspace> \
  -team-id=<team-id> -access=write
```

## Common Scenarios

### Scenario 1: Simple AWS S3 + CloudFront Static Site

```bash
# 1. Create project
hcptf project create -org=my-org -name=static-sites \
  -description="Static website infrastructure"

# 2. Write config (in Git repo):
# - S3 bucket for hosting
# - CloudFront distribution
# - ACM certificate
# - Route53 DNS record

# 3. Create workspace
hcptf workspace create -org=my-org -name=company-website-prod \
  -terraform-version=1.10.0 \
  -description="Company website production infrastructure"

# 4. Set variables
hcptf variable create -org=my-org -workspace=company-website-prod \
  -key=domain_name -value=www.example.com

hcptf variable create -org=my-org -workspace=company-website-prod \
  -key=environment -value=production

# AWS credentials
hcptf variable create -org=my-org -workspace=company-website-prod \
  -key=AWS_ACCESS_KEY_ID -value=$AWS_KEY -category=env -sensitive

hcptf variable create -org=my-org -workspace=company-website-prod \
  -key=AWS_SECRET_ACCESS_KEY -value=$AWS_SECRET -category=env -sensitive

# 5. Deploy
hcptf workspace run create -org=my-org -name=company-website-prod \
  -message="Initial website infrastructure"

# 6. Review and apply
RUN_ID=$(hcptf my-org company-website-prod -output=json | jq -r '.CurrentRunID')
hcptf my-org company-website-prod runs $RUN_ID show
hcptf my-org company-website-prod runs $RUN_ID apply \
  -comment="Approved: initial website setup"

# 7. Get outputs
hcptf state outputs -org=my-org -workspace=company-website-prod
# Shows: cloudfront_domain, bucket_name, etc.
```

### Scenario 2: Multi-Environment Application with Shared Networking

```bash
# 1. Create project
hcptf project create -org=my-org -name=app-platform \
  -description="Application platform infrastructure"

# 2. Create networking workspace (shared)
hcptf workspace create -org=my-org -name=networking-shared \
  -terraform-version=1.10.0 \
  -description="Shared VPC and networking"

hcptf variable create -org=my-org -workspace=networking-shared \
  -key=vpc_cidr -value="10.0.0.0/16"

hcptf variable create -org=my-org -workspace=networking-shared \
  -key=AWS_ACCESS_KEY_ID -value=$AWS_KEY -category=env -sensitive

hcptf variable create -org=my-org -workspace=networking-shared \
  -key=AWS_SECRET_ACCESS_KEY -value=$AWS_SECRET -category=env -sensitive

# Deploy networking first
hcptf workspace run create -org=my-org -name=networking-shared \
  -message="Create shared VPC"

RUN_ID=$(hcptf my-org networking-shared -output=json | jq -r '.CurrentRunID')
hcptf my-org networking-shared runs $RUN_ID apply -comment="Approved"

# Get networking outputs
hcptf state outputs -org=my-org -workspace=networking-shared
# Note: vpc_id, subnet_ids for downstream workspaces

# 3. Create application workspace (per environment)
for env in staging prod; do
  hcptf workspace create -org=my-org -name=api-service-$env \
    -terraform-version=1.10.0 \
    -description="API service - $env"

  hcptf variable create -org=my-org -workspace=api-service-$env \
    -key=environment -value=$env

  hcptf variable create -org=my-org -workspace=api-service-$env \
    -key=AWS_ACCESS_KEY_ID -value=$AWS_KEY -category=env -sensitive

  hcptf variable create -org=my-org -workspace=api-service-$env \
    -key=AWS_SECRET_ACCESS_KEY -value=$AWS_SECRET -category=env -sensitive
done

# Set environment-specific variables
hcptf variable create -org=my-org -workspace=api-service-staging \
  -key=instance_type -value=t3.small

hcptf variable create -org=my-org -workspace=api-service-prod \
  -key=instance_type -value=t3.large

# 4. Deploy staging first
hcptf workspace run create -org=my-org -name=api-service-staging \
  -message="Initial staging deployment"

RUN_ID=$(hcptf my-org api-service-staging -output=json | jq -r '.CurrentRunID')
hcptf my-org api-service-staging runs $RUN_ID apply -comment="Staging approved"

# 5. Then deploy production
hcptf workspace run create -org=my-org -name=api-service-prod \
  -message="Initial production deployment"

RUN_ID=$(hcptf my-org api-service-prod -output=json | jq -r '.CurrentRunID')
hcptf my-org api-service-prod runs $RUN_ID show
# Review plan carefully before applying production
hcptf my-org api-service-prod runs $RUN_ID apply \
  -comment="Production deployment approved after staging validation"
```

### Scenario 3: Quick Single-Workspace Deployment

The simplest possible greenfield setup for a single workspace:

```bash
# 1. Create workspace
hcptf workspace create -org=my-org -name=my-app-prod \
  -terraform-version=1.10.0

# 2. Set all variables at once
hcptf variable create -org=my-org -workspace=my-app-prod \
  -key=environment -value=prod
hcptf variable create -org=my-org -workspace=my-app-prod \
  -key=AWS_ACCESS_KEY_ID -value=$AWS_KEY -category=env -sensitive
hcptf variable create -org=my-org -workspace=my-app-prod \
  -key=AWS_SECRET_ACCESS_KEY -value=$AWS_SECRET -category=env -sensitive

# 3. Trigger run (assumes VCS is connected or config is uploaded)
hcptf workspace run create -org=my-org -name=my-app-prod \
  -message="Initial deployment"

# 4. Apply
RUN_ID=$(hcptf my-org my-app-prod -output=json | jq -r '.CurrentRunID')
hcptf my-org my-app-prod runs $RUN_ID apply -comment="Ship it"

# 5. Verify
hcptf state outputs -org=my-org -workspace=my-app-prod
```

## Post-Deployment Checklist

After the initial deployment, consider setting up:

```bash
# Enable health assessments for drift detection
hcptf workspace update -org=<org> -name=<workspace> -assessments-enabled=true

# Set up run notifications
hcptf notification create -org=<org> -workspace=<workspace> \
  -name=team-alerts -destination-type=slack \
  -url=$SLACK_WEBHOOK -triggers=run:errored,assessment:drifted

# Configure team access
hcptf teamaccess create -org=<org> -workspace=<workspace> \
  -team-id=<team-id> -access=write

# Apply policy sets
hcptf policyset update -organization=<org> -id=<policy-set-id> \
  -workspace=<workspace>

# Verify everything with Explorer
hcptf explorer query -org=<org> -type=workspaces \
  -filter="workspace-name:<workspace>" \
  -fields=workspace-name,workspace-terraform-version,providers,drifted,current-run-status
```

## Tips and Best Practices

1. **Start with a plan**: Design your project/workspace hierarchy before creating anything
2. **Use projects**: Group related workspaces into projects for better organization
3. **Sensitive variables**: Always use `-sensitive` for credentials and secrets
4. **VCS-driven workflows**: Connect workspaces to VCS repos for automated deployments
5. **Pin versions**: Use specific Terraform and provider version constraints
6. **Review before apply**: Disable auto-apply for production workspaces
7. **Enable drift detection**: Turn on health assessments after initial deployment
8. **Set up notifications**: Alert your team on run failures and drift
9. **Use OIDC**: Prefer OIDC-based cloud auth over static credentials
10. **Tag everything**: Use consistent tags/labels across all resources

## Troubleshooting

**Workspace creation fails:**

- Verify organization name: `hcptf organization list`
- Check permissions: `hcptf account show`
- Workspace name must be unique within the org

**Variable creation fails:**

- Workspace must exist first
- Variable keys are case-sensitive
- Use `-category=env` for environment variables (default is `terraform`)

**Run fails with authentication errors:**

- Verify cloud credentials are set as environment variables with `-category=env`
- Check that credential variable keys match provider expectations
- For AWS: `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`
- For Azure: `ARM_CLIENT_ID`, `ARM_CLIENT_SECRET`, `ARM_TENANT_ID`, `ARM_SUBSCRIPTION_ID`
- For GCP: `GOOGLE_CREDENTIALS` or `GOOGLE_PROJECT`

**Run fails with "No configuration uploaded":**

- Workspace needs VCS connection or manual config upload
- Verify VCS connection is working
- Check that Terraform files exist in the repo root (or configured working directory)

**Plan shows unexpected changes:**

- Review variable values: `hcptf variable list -org=<org> -workspace=<workspace>`
- Check Terraform version compatibility
- Verify provider version constraints

## Related Commands

- `hcptf login` - Authenticate with HCP Terraform
- `hcptf account show` - Verify authentication
- `hcptf organization list` - List organizations
- `hcptf project create` - Create a project
- `hcptf project list` - List projects
- `hcptf workspace create` - Create a workspace
- `hcptf workspace read` - View workspace details
- `hcptf workspace update` - Update workspace settings
- `hcptf variable create` - Create a variable
- `hcptf variable list` - List variables
- `hcptf workspace run create` - Create a run
- `hcptf workspace run show` - View run details
- `hcptf workspace run apply` - Apply a run
- `hcptf state outputs` - View state outputs
- `hcptf state list` - List state versions
- `hcptf notification create` - Set up notifications
- `hcptf teamaccess create` - Configure team access
- `hcptf oauthtoken list` - List VCS OAuth tokens
- `hcptf explorer query` - Query resources across org
- `hcptf awsoidc create` - Set up AWS OIDC trust

## Agent Considerations

When building agents that handle greenfield deployments:

1. **Gather requirements first**: Ask about cloud provider, environment strategy, and team structure before creating anything
2. **Validate authentication**: Always verify the token is valid before attempting resource creation
3. **Check for existing resources**: Look for existing projects and workspaces to avoid duplicates
4. **Sensitive data handling**: Never log or display sensitive variable values
5. **Sequential operations**: Create resources in dependency order (org > project > workspace > variables > run)
6. **Wait for runs**: Don't apply immediately; present the plan and ask for confirmation
7. **Provide next steps**: After deployment, suggest drift detection, notifications, and policy setup
