---
name: workspace-to-stack
description: Refactor existing HCP Terraform workspace-based infrastructure into a Terraform Stack. Use when migrating workspaces to stacks for multi-environment orchestration, consolidating related workspaces into a single stack, or breaking monolithic workspaces into stack components with coordinated deployments.
---

# Workspace-to-Stack Refactoring Skill

## Overview

This skill helps refactor existing HCP Terraform workspace-based infrastructure into a Terraform Stack. Stacks enable orchestrating multiple configurations as a cohesive unit with shared inputs, cross-component dependencies, and coordinated multi-environment deployments. This workflow is ideal when a workspace has grown complex, spans multiple environments, or would benefit from stack-level orchestration.

## Prerequisites

- Authenticated with `hcptf` CLI (`TFE_TOKEN` or `~/.terraform.d/credentials.terraform.io`)
- Read access to the source workspace(s)
- Write access to the target project for creating the stack
- VCS repository access for both the existing workspace config and the new stack config
- OAuth token ID for VCS connection (`hcptf oauthtoken list -org=<org>`)

## Core Concepts

**Why Refactor to a Stack?**

- **Multi-environment orchestration**: Deploy the same config to dev, staging, and prod with different inputs
- **Component isolation**: Break a monolithic workspace into smaller, independently deployable components
- **Coordinated deployments**: Deploy related components together with dependency ordering
- **Shared inputs**: Define variables once and pass them to multiple components

**Stack Architecture:**

- **Stack**: Top-level container associated with a VCS repo
- **Components**: Individual Terraform configurations within the stack (defined in `.tfcomponent.hcl`)
- **Deployments**: Instances of the stack deployed to specific environments (defined in `.tfdeploy.hcl`)
- **Stack State**: Per-component, per-deployment state management

**Migration Approach:**

1. Audit the existing workspace to understand its configuration
2. Design the stack structure (components, deployments, inputs)
3. Create the stack configuration files in a new or existing repo
4. Create the stack in HCP Terraform
5. Import or recreate state (depending on migration strategy)
6. Validate and deploy
7. Decommission the old workspace

## Workflow

### 1. Audit the Existing Workspace

**Get workspace details:**

```bash
# View workspace overview
hcptf <org> <workspace>

# Get detailed info including VCS, Terraform version, and run status
hcptf workspace read -org=<org> -name=<workspace> -output=json
```

**Get current configuration source:**

```bash
# Get VCS info from latest run
RUN_ID=$(hcptf <org> <workspace> -output=json | jq -r '.CurrentRunID')
hcptf <org> <workspace> runs $RUN_ID configversion

# Shows: RepoIdentifier, Branch, CommitSHA, CommitURL
```

**Get current variables:**

```bash
# List all workspace variables
hcptf variable list -org=<org> -workspace=<workspace>

# Export as JSON for reference
hcptf variable list -org=<org> -workspace=<workspace> -output=json
```

**Get current state outputs:**

```bash
# View state outputs (these may be consumed by other workspaces)
hcptf state outputs -org=<org> -workspace=<workspace>
```

**Check for run triggers and dependencies:**

```bash
# Check notification configurations
hcptf notification list -org=<org> -workspace=<workspace>

# Check team access
hcptf teamaccess list -org=<org> -workspace=<workspace>
```

**Get resource inventory via Explorer:**

```bash
# See providers and modules used
hcptf explorer query -org=<org> -type=workspaces \
  -filter="workspace-name:<workspace>" \
  -fields=workspace-name,providers,modules,workspace-terraform-version
```

### 2. Design the Stack Structure

Based on the audit, decide how to organize the stack:

**Single-component stack** (simplest migration):

- One component wrapping the entire existing config
- Multiple deployments for environments (dev, staging, prod)
- Good when the workspace is already well-structured

**Multi-component stack** (recommended for complex workspaces):

- Split networking, compute, database, etc. into separate components
- Define dependencies between components
- Each component can be deployed independently

**Decision Matrix:**

| Current State                         | Recommended Approach                             |
| ------------------------------------- | ------------------------------------------------ |
| Single workspace, single environment  | Single component, add deployments for more envs  |
| Single workspace, large config        | Split into multiple components by resource group |
| Multiple workspaces sharing config    | One component per workspace, shared inputs       |
| Multiple workspaces with run triggers | Components with explicit dependencies            |

### 3. Create Stack Configuration Files

**Repository structure for a stack:**

```
my-stack-repo/
  components/
    networking/
      main.tf
      variables.tf
      outputs.tf
    compute/
      main.tf
      variables.tf
      outputs.tf
  deployments.tfdeploy.hcl
  components.tfcomponent.hcl
```

**Define components (components.tfcomponent.hcl):**

```hcl
# Single-component example (wrapping existing workspace config)
component "app" {
  source = "./components/app"

  inputs = {
    environment = var.environment
    region      = var.region
    instance_type = var.instance_type
  }

  providers = {
    aws = provider.aws.this
  }
}

# Multi-component example
component "networking" {
  source = "./components/networking"

  inputs = {
    environment = var.environment
    region      = var.region
    cidr_block  = var.vpc_cidr
  }

  providers = {
    aws = provider.aws.this
  }
}

component "compute" {
  source = "./components/compute"

  inputs = {
    environment   = var.environment
    vpc_id        = component.networking.vpc_id
    subnet_ids    = component.networking.subnet_ids
    instance_type = var.instance_type
  }

  providers = {
    aws = provider.aws.this
  }
}
```

**Define deployments (deployments.tfdeploy.hcl):**

```hcl
identity_token "aws" {
  audience = ["terraform-stacks-private-preview"]
}

deployment "development" {
  inputs = {
    environment   = "dev"
    region        = "us-east-1"
    instance_type = "t3.micro"
  }
}

deployment "staging" {
  inputs = {
    environment   = "staging"
    region        = "us-east-1"
    instance_type = "t3.small"
  }
}

deployment "production" {
  inputs = {
    environment   = "prod"
    region        = "us-east-1"
    instance_type = "t3.large"
  }
}
```

### 4. Create the Stack in HCP Terraform

**Find or create a project:**

```bash
# List existing projects
hcptf project list -org=<org>

# Create a new project if needed
hcptf project create -org=<org> -name=<project-name> \
  -description="Migrated from workspace <workspace>"

# Get project ID
PROJECT_ID=$(hcptf project list -org=<org> -output=json | \
  jq -r '.[] | select(.Name=="<project-name>") | .ID')
```

**Get OAuth token for VCS:**

```bash
# List OAuth tokens
hcptf oauthtoken list -org=<org>

# Note the OAuth token ID for your VCS provider
```

**Create the stack:**

```bash
hcptf stack create \
  -name=<stack-name> \
  -project-id=$PROJECT_ID \
  -vcs-identifier=<org>/<repo> \
  -vcs-branch=main \
  -oauth-token-id=<oauth-token-id> \
  -description="Migrated from workspace <workspace>"

# Note the stack ID from output
```

### 5. Trigger Initial Deployment

**Fetch configuration from VCS:**

```bash
# Trigger stack to pull config from VCS
hcptf stack deployment create -stack-id=<stack-id>
```

**Monitor configuration status:**

```bash
# List configurations
hcptf stack configuration list -stack-id=<stack-id>

# View configuration details
hcptf stack configuration read -id=<config-id>
```

**Monitor deployments:**

```bash
# List deployments
hcptf stack deployment list -stack-id=<stack-id>

# View deployment details
hcptf stack deployment read -deployment-id=<deployment-id>
```

### 6. Validate the Stack

**Check stack state:**

```bash
# List state versions
hcptf stack state list -stack-id=<stack-id>

# View state details
hcptf stack state read -id=<state-id>
```

**Compare outputs:**

```bash
# Compare stack outputs with original workspace outputs
# Original:
hcptf state outputs -org=<org> -workspace=<old-workspace>

# New stack:
hcptf stack state read -id=<state-id>
```

### 7. Decommission the Old Workspace

Only after validating the stack is working correctly:

```bash
# Lock the old workspace to prevent further runs
hcptf workspace update -org=<org> -name=<old-workspace> -locked=true

# Add a description noting the migration
hcptf workspace update -org=<org> -name=<old-workspace> \
  -description="DEPRECATED: Migrated to stack <stack-name>. Do not use."

# Eventually delete (use with caution!)
# hcptf workspace delete -org=<org> -name=<old-workspace> -force
```

## Common Scenarios

### Scenario 1: Simple Workspace to Single-Component Stack

A single workspace managing a web application. Goal: enable multi-environment deployments.

```bash
# 1. Audit existing workspace
hcptf my-org web-app
# Shows: TerraformVersion 1.10.0, VCS connected, 12 resources

hcptf variable list -org=my-org -workspace=web-app
# Shows: environment=prod, region=us-east-1, instance_type=t3.large

RUN_ID=$(hcptf my-org web-app -output=json | jq -r '.CurrentRunID')
hcptf my-org web-app runs $RUN_ID configversion
# Shows: RepoIdentifier: my-org/web-app-infra, Branch: main

# 2. Create stack config files in repo (see step 3 above)
# - Copy existing .tf files to components/app/
# - Create components.tfcomponent.hcl
# - Create deployments.tfdeploy.hcl with dev/staging/prod

# 3. Create stack
PROJECT_ID=$(hcptf project list -org=my-org -output=json | \
  jq -r '.[] | select(.Name=="web-platform") | .ID')

hcptf stack create -name=web-app-stack -project-id=$PROJECT_ID \
  -vcs-identifier=my-org/web-app-infra \
  -vcs-branch=main \
  -oauth-token-id=ot-abc123 \
  -description="Web application stack (migrated from web-app workspace)"

# 4. Deploy
STACK_ID="stk-xyz789"  # from create output
hcptf stack deployment create -stack-id=$STACK_ID

# 5. Monitor
hcptf stack configuration list -stack-id=$STACK_ID
hcptf stack deployment list -stack-id=$STACK_ID

# 6. Lock old workspace after validation
hcptf workspace update -org=my-org -name=web-app -locked=true \
  -description="DEPRECATED: Migrated to stack web-app-stack"
```

### Scenario 2: Multiple Related Workspaces to Multi-Component Stack

Three workspaces (networking, compute, database) that share a VPC. Goal: unify into one stack.

```bash
# 1. Audit all workspaces
for ws in networking compute database; do
  echo "=== $ws ==="
  hcptf my-org $ws
  hcptf variable list -org=my-org -workspace=$ws
  hcptf state outputs -org=my-org -workspace=$ws
done

# 2. Document cross-workspace dependencies
# networking outputs: vpc_id, subnet_ids, security_group_ids
# compute uses: vpc_id, subnet_ids from networking
# database uses: vpc_id, subnet_ids from networking

# 3. Create stack structure:
#    components/networking/ - from networking workspace config
#    components/compute/    - from compute workspace config
#    components/database/   - from database workspace config
#    components.tfcomponent.hcl - define component dependencies
#    deployments.tfdeploy.hcl  - define environments

# 4. Create and deploy stack
hcptf stack create -name=platform-stack -project-id=$PROJECT_ID \
  -vcs-identifier=my-org/platform-infra \
  -oauth-token-id=ot-abc123

STACK_ID="stk-abc123"
hcptf stack deployment create -stack-id=$STACK_ID

# 5. Monitor all components
hcptf stack configuration list -stack-id=$STACK_ID
hcptf stack deployment list -stack-id=$STACK_ID

# 6. Lock old workspaces after validation
for ws in networking compute database; do
  hcptf workspace update -org=my-org -name=$ws -locked=true \
    -description="DEPRECATED: Migrated to stack platform-stack"
done
```

### Scenario 3: Workspace with Multiple Environments to Stack Deployments

A workspace that's duplicated for each environment (app-dev, app-staging, app-prod). Goal: consolidate into a single stack with deployments.

```bash
# 1. Audit all environment workspaces
hcptf explorer query -org=my-org -type=workspaces \
  -filter="workspace-name:app-" \
  -fields=workspace-name,workspace-terraform-version,providers,current-run-status

# 2. Compare variables across environments
for env in dev staging prod; do
  echo "=== app-$env ==="
  hcptf variable list -org=my-org -workspace=app-$env -output=json
done

# 3. Identify what varies between environments (these become deployment inputs)
# Common differences: instance_type, replica_count, domain_name, etc.

# 4. Create stack with one component and three deployments
# - components/app/ contains the shared Terraform config
# - deployments.tfdeploy.hcl defines dev, staging, prod with different inputs

# 5. Create stack
hcptf stack create -name=app-stack -project-id=$PROJECT_ID \
  -vcs-identifier=my-org/app-infra \
  -oauth-token-id=ot-abc123

# 6. Deploy and validate each environment
hcptf stack deployment create -stack-id=$STACK_ID

# 7. Decommission old per-environment workspaces
for env in dev staging prod; do
  hcptf workspace update -org=my-org -name=app-$env -locked=true \
    -description="DEPRECATED: Migrated to stack app-stack, deployment $env"
done
```

## Tips and Best Practices

1. **Migrate incrementally**: Start with one environment (dev) to validate the stack, then add staging/prod
2. **Keep the old workspace**: Lock it rather than deleting until the stack is fully validated
3. **Document variable mapping**: Create a clear mapping of workspace variables to stack deployment inputs
4. **Test with speculative plans**: Use `-speculative-enabled` on the stack to test PRs
5. **Preserve state carefully**: If importing existing state, plan thoroughly to avoid resource recreation
6. **Update downstream consumers**: If other workspaces depend on outputs from the migrated workspace, update them to use stack outputs
7. **Use the same Terraform version**: Match the stack's Terraform version to the workspace's version initially
8. **Review provider authentication**: Stacks may use different auth mechanisms (OIDC tokens) than workspaces

## Troubleshooting

**Stack creation fails with VCS error:**

- Verify OAuth token ID is correct: `hcptf oauthtoken list -org=<org>`
- Ensure repo identifier format is correct (`org/repo`)
- Check that the VCS branch exists

**Configuration fetch fails:**

- Verify stack configuration files exist in the repo root
- Check `.tfcomponent.hcl` and `.tfdeploy.hcl` syntax
- Ensure the VCS connection has read access to the repo

**Deployment fails with provider errors:**

- Stacks use provider configurations differently
- Check that provider blocks are properly referenced in components
- Verify OIDC or credential setup for the deployment context

**State import issues:**

- Stack state is per-component, per-deployment
- Cannot directly import workspace state into a stack
- May need to recreate resources or use `terraform import` within each component

**Old workspace still receiving VCS triggers:**

- Disconnect VCS from the old workspace before locking
- Or delete the VCS connection: update workspace to remove VCS settings

## Related Commands

- `hcptf workspace read` - View workspace details
- `hcptf variable list` - List workspace variables
- `hcptf state outputs` - View workspace state outputs
- `hcptf configversion read` - Get VCS info from workspace
- `hcptf project create` - Create a project for the stack
- `hcptf project list` - List projects to find project ID
- `hcptf oauthtoken list` - List OAuth tokens for VCS
- `hcptf stack create` - Create a new stack
- `hcptf stack configuration list` - List stack configurations
- `hcptf stack deployment create` - Trigger a stack deployment
- `hcptf stack deployment list` - List deployments
- `hcptf stack state list` - List stack state versions
- `hcptf workspace update` - Lock/update old workspace
- `hcptf explorer query` - Query workspaces and resources across org

## Agent Considerations

When building agents that handle workspace-to-stack migration:

1. **Full audit first**: Always gather complete workspace information before suggesting a stack design
2. **Get approval for design**: Present the proposed component/deployment structure before creating anything
3. **Warn about state**: Clearly communicate that state cannot be directly migrated; resources may need recreation
4. **Preserve existing infrastructure**: Never delete or modify the running workspace until the stack is validated
5. **Track dependencies**: Identify and document all cross-workspace dependencies before migration
6. **Suggest incremental approach**: Recommend migrating one environment at a time
