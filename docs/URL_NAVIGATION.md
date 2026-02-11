# URL-Style Navigation Guide

The `hcptf` CLI supports an intuitive URL-style navigation pattern alongside traditional command syntax. This allows you to explore your HCP Terraform resources using a path-like approach.

## Concept

Instead of using verbose flag-based commands, you can navigate resources like browsing a filesystem or URL:

```
hcptf <org> [<workspace>] [<resource>] [<action>]
```

## Examples

### Organization Level

**Show organization details:**
```bash
hcptf my-org
```
Equivalent to: `hcptf organization show -org=my-org`

**List resources in organization:**
```bash
hcptf my-org workspaces      # List all workspaces
hcptf my-org projects         # List all projects
hcptf my-org teams            # List all teams
hcptf my-org policies         # List all policies
hcptf my-org policysets       # List all policy sets
```

Equivalent to:
```bash
hcptf workspace list -org=my-org
hcptf project list -org=my-org
hcptf team list -org=my-org
hcptf policy list -org=my-org
hcptf policyset list -org=my-org
```

### Workspace Level

**Show workspace details:**
```bash
hcptf my-org my-workspace
```
Equivalent to: `hcptf workspace read -org=my-org -workspace=my-workspace`

**List workspace resources:**
```bash
hcptf my-org my-workspace runs         # List runs
hcptf my-org my-workspace variables    # List variables
hcptf my-org my-workspace state        # List state versions
```

Equivalent to:
```bash
hcptf run list -org=my-org -workspace=my-workspace
hcptf variable list -org=my-org -workspace=my-workspace
hcptf state list -org=my-org -workspace=my-workspace
```

### Run Operations

**List runs:**
```bash
hcptf my-org my-workspace runs
hcptf my-org my-workspace runs list    # Explicit list
```

**Show run details:**
```bash
hcptf my-org my-workspace runs run-abc123
```
Equivalent to: `hcptf run show -id=run-abc123`

**Apply a run:**
```bash
hcptf my-org my-workspace runs run-abc123 apply
```
Equivalent to: `hcptf run apply -id=run-abc123`

**Other run actions:**
```bash
hcptf my-org my-workspace runs run-abc123 discard
hcptf my-org my-workspace runs run-abc123 cancel
```

### State Operations

**List state versions:**
```bash
hcptf my-org my-workspace state
hcptf my-org my-workspace state list    # Explicit list
```

**Show state outputs:**
```bash
hcptf my-org my-workspace state outputs
```
Equivalent to: `hcptf state outputs -org=my-org -workspace=my-workspace`

## Traditional Command Syntax

The traditional command syntax is still fully supported and preferred for:
- Automation and scripting (more explicit)
- Complex operations with many flags
- Operations that don't fit the path model

Both styles can be used interchangeably based on your preference!

## When URL-Style is Applied

The URL-style routing is automatically detected when:
1. The first argument doesn't start with `-` (not a flag)
2. The first argument is not a known command name

If either condition fails, the CLI falls back to traditional command parsing.

## Limitations

URL-style navigation currently supports:
- ✅ Organization operations
- ✅ Workspace operations
- ✅ Run operations
- ✅ Variable listing
- ✅ State operations
- ✅ Resource listings (workspaces, projects, teams, policies, policysets)

Not yet supported via URL-style:
- ❌ Create operations (use traditional syntax: `hcptf workspace create ...`)
- ❌ Update operations (use traditional syntax: `hcptf workspace update ...`)
- ❌ Delete operations (use traditional syntax: `hcptf workspace delete ...`)
- ❌ Advanced operations with multiple required parameters

For these operations, use the traditional command syntax with explicit flags.

## Examples Side-by-Side

| URL-Style | Traditional |
|-----------|-------------|
| `hcptf my-org` | `hcptf organization show -org=my-org` |
| `hcptf my-org workspaces` | `hcptf workspace list -org=my-org` |
| `hcptf my-org my-ws` | `hcptf workspace read -org=my-org -workspace=my-ws` |
| `hcptf my-org my-ws runs` | `hcptf run list -org=my-org -workspace=my-ws` |
| `hcptf my-org my-ws runs run-123` | `hcptf run show -id=run-123` |
| `hcptf my-org my-ws runs run-123 apply` | `hcptf run apply -id=run-123` |
| `hcptf my-org my-ws variables` | `hcptf variable list -org=my-org -workspace=my-ws` |
| `hcptf my-org my-ws state outputs` | `hcptf state outputs -org=my-org -workspace=my-ws` |
| `hcptf my-org teams` | `hcptf team list -org=my-org` |

## Tips

1. **Tab completion**: URL-style navigation works great with shell tab completion
2. **Exploration**: Use this style when exploring resources interactively
3. **Scripting**: Use traditional syntax for scripts and CI/CD (more explicit)
4. **Mix and match**: You can use both styles in the same workflow
