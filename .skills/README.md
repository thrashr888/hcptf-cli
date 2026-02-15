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
