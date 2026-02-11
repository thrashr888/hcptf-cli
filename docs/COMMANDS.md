# Command Reference

229 commands across 59 resource types. All commands support `-output=json` for machine-readable output.

## Authentication (2)

| Command | Description |
|---------|-------------|
| `login` | Authenticate and save credentials |
| `logout` | Remove saved credentials |

## Account (3)

| Command | Description |
|---------|-------------|
| `account create` | Create a new HCP Terraform account |
| `account show` | Show current account details |
| `account update` | Update account email or username |

## Workspace (5)

| Command | Description |
|---------|-------------|
| `workspace list` | List workspaces in an organization |
| `workspace create` | Create a new workspace |
| `workspace read` | Show workspace details |
| `workspace update` | Update workspace settings |
| `workspace delete` | Delete a workspace |

## Run (6)

| Command | Description |
|---------|-------------|
| `run list` | List runs for a workspace |
| `run create` | Trigger a new run |
| `run show` | Show run details |
| `run apply` | Approve and apply a run |
| `run discard` | Discard a run |
| `run cancel` | Cancel a running run |

## Organization (5)

| Command | Description |
|---------|-------------|
| `organization list` | List accessible organizations |
| `organization create` | Create a new organization |
| `organization show` | Show organization details |
| `organization update` | Update organization settings |
| `organization delete` | Delete an organization |

## Variable (4)

| Command | Description |
|---------|-------------|
| `variable list` | List variables for a workspace |
| `variable create` | Create a new variable |
| `variable update` | Update a variable |
| `variable delete` | Delete a variable |

## Team (6)

| Command | Description |
|---------|-------------|
| `team list` | List teams in an organization |
| `team create` | Create a new team |
| `team show` | Show team details |
| `team delete` | Delete a team |
| `team add-member` | Add a user to a team |
| `team remove-member` | Remove a user from a team |

## Project (5)

| Command | Description |
|---------|-------------|
| `project list` | List projects in an organization |
| `project create` | Create a new project |
| `project read` | Show project details |
| `project update` | Update project settings |
| `project delete` | Delete a project |

## State (3)

| Command | Description |
|---------|-------------|
| `state list` | List state versions for a workspace |
| `state read` | Read state version details |
| `state outputs` | Display outputs from current state |

## Policy (5)

| Command | Description |
|---------|-------------|
| `policy list` | List policies in an organization |
| `policy create` | Create a policy with content upload |
| `policy read` | Read policy details and content |
| `policy update` | Update policy settings and content |
| `policy delete` | Delete a policy |

## Policy Set (7)

| Command | Description |
|---------|-------------|
| `policyset list` | List policy sets |
| `policyset create` | Create a new policy set |
| `policyset read` | Read policy set details |
| `policyset update` | Update policy set settings |
| `policyset delete` | Delete a policy set |
| `policyset add-policy` | Add a policy to a set |
| `policyset remove-policy` | Remove a policy from a set |

## SSH Key (5)

| Command | Description |
|---------|-------------|
| `sshkey list` | List SSH keys in an organization |
| `sshkey create` | Create a new SSH key |
| `sshkey read` | Read SSH key details |
| `sshkey update` | Update SSH key name |
| `sshkey delete` | Delete an SSH key |

## Notification (6)

| Command | Description |
|---------|-------------|
| `notification list` | List notification configurations |
| `notification create` | Create a notification (Slack, webhook, Teams, email) |
| `notification read` | Read notification details |
| `notification update` | Update notification settings |
| `notification delete` | Delete a notification |
| `notification verify` | Test notification delivery |

## Variable Set (10)

| Command | Description |
|---------|-------------|
| `variableset list` | List variable sets in an organization |
| `variableset create` | Create a new variable set |
| `variableset read` | Read variable set details |
| `variableset update` | Update variable set settings |
| `variableset delete` | Delete a variable set |
| `variableset apply` | Apply variable set to workspaces/projects |
| `variableset variable list` | List variables in a variable set |
| `variableset variable create` | Create variable in a variable set |
| `variableset variable update` | Update variable in a variable set |
| `variableset variable delete` | Delete variable from a variable set |

## Agent Pool (8)

| Command | Description |
|---------|-------------|
| `agentpool list` | List agent pools in an organization |
| `agentpool create` | Create a new agent pool |
| `agentpool read` | Read agent pool details |
| `agentpool update` | Update agent pool settings |
| `agentpool delete` | Delete an agent pool |
| `agentpool token-create` | Create agent authentication token |
| `agentpool token-list` | List agent tokens |
| `agentpool token-delete` | Delete an agent token |

## Run Task (7)

| Command | Description |
|---------|-------------|
| `runtask list` | List run tasks in an organization |
| `runtask create` | Create a new run task |
| `runtask read` | Read run task details |
| `runtask update` | Update run task settings |
| `runtask delete` | Delete a run task |
| `runtask attach` | Attach run task to a workspace |
| `runtask detach` | Detach run task from a workspace |

## OAuth Client (5)

| Command | Description |
|---------|-------------|
| `oauthclient list` | List OAuth clients in an organization |
| `oauthclient create` | Create VCS OAuth client (GitHub, GitLab, etc.) |
| `oauthclient read` | Read OAuth client details |
| `oauthclient update` | Update OAuth client settings |
| `oauthclient delete` | Delete an OAuth client |

## OAuth Token (3)

| Command | Description |
|---------|-------------|
| `oauthtoken list` | List OAuth tokens |
| `oauthtoken read` | Read OAuth token details |
| `oauthtoken update` | Update OAuth token SSH key |

## Run Trigger (4)

| Command | Description |
|---------|-------------|
| `runtrigger list` | List run triggers for a workspace |
| `runtrigger create` | Create run trigger to link workspaces |
| `runtrigger read` | Read run trigger details |
| `runtrigger delete` | Delete a run trigger |

## Plan (2)

| Command | Description |
|---------|-------------|
| `plan read` | Read plan details and resource changes |
| `plan logs` | Get plan execution logs |

## Apply (2)

| Command | Description |
|---------|-------------|
| `apply read` | Read apply details and resource changes |
| `apply logs` | Get apply execution logs |

## Configuration Version (4)

| Command | Description |
|---------|-------------|
| `configversion list` | List configuration versions for a workspace |
| `configversion create` | Create new configuration version |
| `configversion read` | Read configuration version details |
| `configversion upload` | Upload configuration files |

## Team Access (5)

| Command | Description |
|---------|-------------|
| `teamaccess list` | List team access for a workspace |
| `teamaccess create` | Grant team access to workspace |
| `teamaccess read` | Read team access details |
| `teamaccess update` | Update team workspace permissions |
| `teamaccess delete` | Remove team access from workspace |

## Project Team Access (5)

| Command | Description |
|---------|-------------|
| `projectteamaccess list` | List team access for a project |
| `projectteamaccess create` | Grant team access to project |
| `projectteamaccess read` | Read project team access details |
| `projectteamaccess update` | Update team project permissions |
| `projectteamaccess delete` | Remove team access from project |

## Registry Module (6)

| Command | Description |
|---------|-------------|
| `registrymodule list` | List private registry modules |
| `registrymodule create` | Create/publish a new module |
| `registrymodule read` | Read module details |
| `registrymodule delete` | Delete a module |
| `registrymodule create-version` | Publish new module version |
| `registrymodule delete-version` | Delete a module version |

## Registry Provider (4)

| Command | Description |
|---------|-------------|
| `registryprovider list` | List private registry providers |
| `registryprovider create` | Create a new provider |
| `registryprovider read` | Read provider details |
| `registryprovider delete` | Delete a provider |

## Registry Provider Version (3)

| Command | Description |
|---------|-------------|
| `registryproviderversion create` | Create provider version |
| `registryproviderversion read` | Read provider version details |
| `registryproviderversion delete` | Delete provider version |

## Registry Provider Platform (3)

| Command | Description |
|---------|-------------|
| `registryproviderplatform create` | Add platform (OS/Arch) for provider |
| `registryproviderplatform read` | Read platform details |
| `registryproviderplatform delete` | Delete platform |

## GPG Key (5)

| Command | Description |
|---------|-------------|
| `gpgkey list` | List GPG keys for provider signing |
| `gpgkey create` | Upload GPG public key |
| `gpgkey read` | Read GPG key details |
| `gpgkey update` | Update GPG key namespace |
| `gpgkey delete` | Delete a GPG key |

## Stack (5)

| Command | Description |
|---------|-------------|
| `stack list` | List Terraform Stacks |
| `stack create` | Create a new stack |
| `stack read` | Read stack details |
| `stack update` | Update stack settings |
| `stack delete` | Delete a stack |

## Stack Configuration (5)

| Command | Description |
|---------|-------------|
| `stackconfiguration list` | List stack configurations |
| `stackconfiguration create` | Create stack configuration |
| `stackconfiguration read` | Read configuration details |
| `stackconfiguration update` | Update configuration |
| `stackconfiguration delete` | Delete configuration |

## Stack Deployment (3)

| Command | Description |
|---------|-------------|
| `stackdeployment list` | List stack deployments |
| `stackdeployment create` | Trigger new deployment |
| `stackdeployment read` | Read deployment details |

## Stack State (2)

| Command | Description |
|---------|-------------|
| `stackstate list` | List stack state versions |
| `stackstate read` | Read stack state details |

## Audit Trail (2)

| Command | Description |
|---------|-------------|
| `audittrail list` | List audit trail events |
| `audittrail read` | Read audit event details |

## Audit Trail Token (4)

| Command | Description |
|---------|-------------|
| `audittrailtoken list` | List audit trail tokens |
| `audittrailtoken create` | Create audit trail token |
| `audittrailtoken read` | Read token details |
| `audittrailtoken delete` | Delete audit trail token |

## System (1)

| Command | Description |
|---------|-------------|
| `version` | Display CLI version |

## Common Patterns

```bash
# List resources
hcptf <resource> list -org=<org>

# Create resource
hcptf <resource> create -org=<org> -name=<name> [options]

# Read resource
hcptf <resource> read -id=<id>

# Update resource
hcptf <resource> update -id=<id> [options]

# Delete resource (with confirmation)
hcptf <resource> delete -id=<id> [-force]
```

Use `hcptf <command> --help` for flags and examples on any command.
