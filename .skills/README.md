# Agent Skills

This directory contains [Agent Skills](https://agentskills.io/) that help AI agents use the hcptf CLI effectively.

## What are Agent Skills?

Agent Skills are folders of instructions, scripts, and resources that AI agents can discover and use to perform tasks more accurately and efficiently. They're supported by leading AI development tools including Claude Code, Cursor, GitHub Copilot, and many others.

## Available Skills

### hcptf-cli

Comprehensive guide for managing HCP Terraform resources using the hcptf command-line tool.

**Use when:**
- Managing HCP Terraform workspaces, runs, or organizations
- Working with Terraform Stacks or private registry resources
- Automating infrastructure deployments
- Querying resource state across organizations

**Key topics covered:**
- Authentication and setup
- Hierarchical command structure (registry, stack, etc.)
- URL-style navigation
- Common workflows (workspace creation, deployments, registry management)
- Output formats (table, JSON)
- Best practices and troubleshooting

### drift

Guide for investigating and resolving infrastructure drift in HCP Terraform workspaces.

**Use when:**
- Infrastructure has changed outside of Terraform
- Investigating drift detection alerts
- Deciding whether to update code or fix infrastructure
- Resolving continuous validation check failures

**Key topics covered:**
- Finding drifted workspaces with Explorer API
- Viewing detailed drift information (what changed, before/after values)
- Getting VCS commit information for the configuration
- Decision matrix for drift resolution strategies
- Common drift scenarios and solutions
- Verifying drift resolution

### version-upgrades

Guide for upgrading Terraform, provider, module, and policy versions in workspaces.

**Use when:**
- Upgrading workspace Terraform version
- Updating provider versions for security patches or new features
- Upgrading modules to latest versions
- Updating policy sets to newer versions
- Planning organization-wide version updates
- Rolling back failed upgrades

**Key topics covered:**
- Finding current versions across workspaces (Explorer API)
- Checking for outdated versions
- Upgrading workspace Terraform version (workspace setting)
- Updating provider/module/policy versions in code (VCS workflow)
- Testing upgrades with speculative runs
- Handling breaking changes and rollbacks
- Bulk upgrade strategies for multiple workspaces

### policy-compliance

Guide for investigating and resolving policy check failures.

**Use when:**
- Runs are blocked by failed policy checks
- Understanding why policies failed
- Deciding whether to fix code or override policies
- Finding policy failures across the organization
- Testing new policies before enforcement
- Troubleshooting policy check issues

**Key topics covered:**
- Detecting policy failures in runs
- Viewing policy check details and results
- Understanding what policies check (using public registry)
- Identifying which resources violated policies
- Decision matrix for remediation (fix code, override, adjust policy)
- Common policy scenarios (CIS benchmarks, tagging, security groups)
- Policy troubleshooting and best practices
- Tracking compliance metrics across organization

### workspace-to-stack

Guide for refactoring existing workspace-based infrastructure into a Terraform Stack.

**Use when:**
- Migrating a workspace to a stack for multi-environment orchestration
- Consolidating multiple related workspaces into a single stack
- Breaking a monolithic workspace into stack components
- Combining per-environment workspaces (app-dev, app-staging, app-prod) into stack deployments

**Key topics covered:**
- Auditing existing workspaces (config, variables, state, dependencies)
- Designing stack structure (single vs multi-component, deployments)
- Creating stack configuration files (`.tfcomponent.hcl`, `.tfdeploy.hcl`)
- Creating and deploying the stack in HCP Terraform
- Validating the migration and comparing outputs
- Decommissioning old workspaces safely
- Common migration scenarios and troubleshooting

### greenfield-deploy

Guide for setting up a brand-new infrastructure project from scratch with HCP Terraform.

**Use when:**
- Starting a new infrastructure project from zero
- Setting up HCP Terraform for the first time
- Creating workspaces, variables, and triggering first deployments
- Onboarding a new team to HCP Terraform
- Building a multi-environment deployment from scratch

**Key topics covered:**
- Planning project and workspace hierarchy
- Verifying authentication and organization setup
- Creating projects and workspaces
- Writing initial Terraform configuration
- Configuring variables and cloud provider credentials
- Triggering and reviewing the first run
- Setting up notifications, team access, and health assessments
- Post-deployment checklist and best practices

## Using Skills

### In Compatible Agents

Agent Skills are automatically discovered by compatible agents when present in the `.skills/` directory. Agents will:

1. Load skill metadata at startup
2. Activate relevant skills based on task context
3. Follow skill instructions and examples
4. Reference additional resources as needed

### Manually

You can also reference skills directly when working with AI assistants:

```
Please use the hcptf-cli skill to help me create a new workspace and deploy infrastructure.
```

## Format

Each skill follows the [Agent Skills specification](https://agentskills.io/specification):

- **SKILL.md**: Required file with YAML frontmatter and Markdown instructions
- **scripts/**: Optional executable code
- **references/**: Optional detailed documentation
- **assets/**: Optional templates and data files

## Contributing

To add or improve skills:

1. Follow the [Agent Skills specification](https://agentskills.io/specification)
2. Keep SKILL.md focused and under 500 lines
3. Include clear examples and common workflows
4. Test with multiple AI agents when possible

## Learn More

- [Agent Skills Website](https://agentskills.io/)
- [Specification](https://agentskills.io/specification)
- [Example Skills](https://github.com/anthropics/skills)
- [Integration Guide](https://agentskills.io/integrate-skills)
