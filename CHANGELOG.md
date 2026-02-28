# Changelog

All notable changes to the HCP Terraform CLI (`hcptf`).

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- **Run show workspace name**: Automatically include workspace relation in `run show` so WorkspaceName is always populated
- **KeyValue output ordering**: Sort keys alphabetically for deterministic, reproducible output
- **Variable list JSON output**: Preserve full variable values in JSON output instead of truncating to 50 characters (truncation still applies to table display)

## [0.5.0] - 2026-02-22

### Added

- **Workspace project reassignment**: `hcptf workspace update -project-id=<id>` to move workspaces between projects
- **Full API coverage for workspace create/update**: Added support for all `WorkspaceCreateOptions` and `WorkspaceUpdateOptions` fields including execution-mode, working-directory, agent-pool-id, VCS repo management, trigger prefixes/patterns, auto-destroy settings, tags, source tracking, and all boolean toggle flags
- **Workspace lock lifecycle**: `workspace lock`, `workspace unlock`, and `workspace force-unlock` commands
- **Run force-execute**: `run force-execute` command to force-execute a pending run
- **Workspace read project metadata**: `workspace read` now includes `ProjectID` and `ProjectName` in output, with automatic project include
- **List/read filter flags**: Added search, tag, status, include, and sort flags to `workspace list`, `run list`, `run show`, and `workspace read`
- **Run org-level listing**: `run list-org` command for listing runs across an organization
- **Policy upload/download**: `policy upload` and `policy download` commands
- **Variable set management**: `variableset remove`, `variableset list-workspace`, `variableset list-project`, `variableset update-workspaces`, `variableset update-stacks` commands with `-query` and `-include` flags on list/read
- **Policy set workspace/project management**: `policyset add-workspace`, `policyset remove-workspace`, `policyset add-workspace-exclusion`, `policyset remove-workspace-exclusion`, `policyset add-project`, `policyset remove-project` commands
- **Expanded policy/policyset flags**: Added `-kind`, `-query`, `-include`, `-overridable`, `-agent-enabled`, `-policy-tool-version`, `-policies-path` flags to policy and policyset commands
- **API coverage Go tests**: Converted `scripts/api-coverage.sh` and `scripts/test-priority.sh` to Go tests in `command/commands_test.go` — now run automatically in CI via `go test ./...`
  - `TestAllCommandsRegistered` — verifies every expected API operation has a registered CLI command
  - `TestAllCommandsHaveHelpAndSynopsis` — asserts all commands have non-empty Help() and Synopsis()
  - `TestAllCommandFilesHaveTests` — reports which command files lack test coverage
  - `TestCommandNamesMatchFiles` — validates command-to-file naming conventions

### Removed

- Deleted `scripts/api-coverage.sh`, `scripts/test-priority.sh`, and `scripts/coverage.sh` — replaced by Go tests or no longer needed

## [0.4.0] - 2026-02-16

### Added

- **State Analyzer Skill**: Comprehensive guide (549 lines) for analyzing Terraform state files
  - Identifies security issues (exposed secrets, public access, unencrypted resources)
  - Finds cost optimization opportunities (over-provisioned instances, unused resources)
  - Detects best practice violations (tagging, naming conventions, deprecated resources)
  - Includes jq-based analysis patterns and example workflows
- **Plan Analyzer Skill**: Guide (476 lines) for reviewing Terraform plans before applying
  - Identifies destructive changes, replacements, and high-risk resources
  - Provides safety checklists and validation scripts
  - Covers common scenarios (routine updates, unexpected changes, database modifications)
- **State Download Command**: `hcptf state download` to retrieve state JSON from HCP Terraform
  - Supports file output (`-output=state.json`) and stdout piping for jq workflows
  - Downloads from current workspace state or specific state version by ID
  - Essential for state analysis and debugging workflows
- **Workspace-to-Stack Skill**: Guide for refactoring existing workspaces into Terraform Stacks
  - Audit, design, migrate, and validate workflows
  - Best practices for component organization and deployment structure
- **Greenfield Deploy Skill**: Guide for setting up new projects from scratch with HCP Terraform
  - Create workspace, configure settings, deploy infrastructure, verify results
- **Dependabot Configuration**: Automated dependency update PRs for Go modules and GitHub Actions
- Updated skill installation instructions with `npx skills add` command

### Changed

- **URL-style navigation is now read-only** for safety (follows kubectl/gh/docker patterns)
  - `hcptf <org> <workspace> runs <run-id> apply` now shows apply details (read-only)
  - To execute apply, use flag-based command: `hcptf run apply -id=<run-id>`
  - This prevents accidental destructive actions when navigating resources
- Simplified command patterns with two-word actions: `apply logs`, `plan logs`
- Updated GitHub Actions to latest versions (Dependabot PR)
- Updated indirect dependencies for bug fixes and compatibility (PR #5)
  - github.com/clipperhouse/uax29/v2: v2.6.0 → v2.7.0 (Unicode text segmentation improvements)
  - github.com/go-test/deep: v1.0.3 → v1.1.1 (Deep equality testing enhancements)
  - github.com/rogpeppe/go-internal: v1.9.0 → v1.14.1 (Internal Go tooling utilities updates)

### Documentation

- Clarified separation between read-only URL-style commands and explicit flag-based action commands
- Added "Actions (Flag-Based Commands)" section to README showing proper usage for destructive operations
- Updated help text to indicate URL-style commands are for navigation only

## [0.3.1] - 2026-02-15

### Security

- **Critical dependency updates**: Updated golang.org/x/crypto from v0.38.0 to v0.48.0 to address security vulnerabilities
  - Fixed unbounded memory consumption in golang.org/x/crypto/ssh
  - Fixed panic vulnerability in golang.org/x/crypto/ssh/agent from malformed messages
- Updated 30+ transitive dependencies to latest stable versions for security and stability

### Changed

- Updated Go version to 1.25
- Updated github.com/hashicorp/go-slug v0.16.8 → v1.0.0
- Updated github.com/hashicorp/jsonapi v1.4.3 → v1.5.0
- Replaced github.com/imdario/mergo with dario.cat/mergo v1.0.2

## [0.3.0] - 2026-02-15

### Changed

- **Information architecture rework**:
  - Canonical nested namespaces are now the primary command model.
  - Removed legacy flattened namespace aliases for normalized command groups (pre-1.0 cleanup).
  - Updated command help and docs to reflect canonical namespace paths.
- **Router model generation**:
  - Replaced hardcoded router known-command and resource keyword fallbacks with command-registry-derived model generation.
  - Added model drift guard tests to ensure router roots stay aligned with the command registry.
- **Command dispatch behavior**:
  - Added implicit GET inference for concise usage (`list`/`read`/`show`) with deterministic ambiguity errors.
  - Standardized delete bypass flags so `-f` and `-y` normalize to `-force` for delete commands.

### Fixed

- Fixed Go formatting issues in command/commands.go, command/login.go, and command/whoami.go to pass CI checks.

## [0.2.0] - 2026-02-14

### Fixed

- Corrected command test helpers to construct mock API clients with hostname-based credentials so new command tests can run reliably in CI without external credentials.
- Added missing changelog entries and follow-up test refactors for command coverage-focused changes.

### Changed

- **BREAKING**: Replaced flat command structure with hierarchical namespaces for registry and stack commands
  - Registry commands now use: `registry module`, `registry provider`, `registry provider version`, `registry provider platform`
  - Stack commands now use: `stack`, `stack configuration`, `stack deployment`, `stack state`
  - Removed legacy commands: `registrymodule`, `registryprovider`, `registryproviderversion`, `registryproviderplatform`, `stackconfiguration`, `stackdeployment`, `stackstate`
  - Benefits: Cleaner command structure, better discoverability, follows HashiCorp CLI best practices

### Added

- **Public Terraform Registry Commands**: Query public registry for providers, modules, and policies
  - `publicregistry provider` - Get provider info (version, description, docs URL)
  - `publicregistry provider versions` - List all available provider versions with protocols and platforms
  - `publicregistry module` - Get module info (version, downloads, verified status)
  - `publicregistry policy` - Get policy set details (included policies, modules, version)
  - `publicregistry policy list` - List all available public Sentinel/OPA policies
  - All commands support table and JSON output formats
  - Similar to [terraform-mcp-server](https://github.com/hashicorp/terraform-mcp-server) capabilities
  - Useful for version upgrade workflows and discovering compliance policies

- **Enhanced Assessment Results**:
  - Added support for Terraform continuous validation checks
  - Display check status (pass/fail/error/unknown) with problem details
  - Parse health-json-redacted endpoint for comprehensive check data
  - Show check instances with per-instance problems
  - Prefer health-json-redacted over json-output for checks
  - URL-style patterns: `hcptf <org> <workspace> assessments` and `hcptf <org> <workspace> runs <run-id> assessment`

- **Configuration Version VCS Info**:
  - Added ingress attributes support to `configversion read`
  - Shows Branch, CommitSHA, CommitURL, CompareURL, RepoIdentifier
  - Essential for drift investigation and code-based remediation workflows
  - Fetches from `/api/v2/configuration-versions/{id}/ingress-attributes` endpoint

- **Agent Skills for Common Workflows**:
  - `.skills/drift/` - Investigate and resolve infrastructure drift
    - Finding drifted workspaces using Explorer API
    - Viewing drift details with assessment results
    - Getting VCS commit information
    - Decision matrix for drift resolution strategies
    - Common drift scenarios (deleted resources, IP changes, tags, cert expiration)
  - `.skills/version-upgrades/` - Upgrade Terraform, provider, module, and policy versions
    - Terraform version upgrades (workspace setting)
    - Provider/module/policy upgrades (VCS workflow with code changes)
    - Using publicregistry commands to find latest versions
    - Complete git clone/edit/commit/push workflows
    - Handling breaking changes and rollbacks
  - `.skills/policy-compliance/` - Investigate and resolve policy check failures
    - Detecting policy failures across workspaces
    - Understanding what policies check (using publicregistry policy)
    - Identifying violating resources
    - Decision matrix for remediation (fix code, override, adjust policy)
    - Common scenarios (CIS benchmarks, tagging, security groups)
    - Tracking compliance metrics across organization
  - `.skills/hcptf-cli/` - Comprehensive CLI usage guide
  - All skills follow [Agent Skills specification](https://agentskills.io/)
  - Automatically discovered by compatible agents (Claude Code, Cursor, GitHub Copilot, etc.)

- Parent commands for hierarchical navigation: `hcptf registry`, `hcptf stack`, and `hcptf publicregistry`
- Explorer API support: `hcptf explorer query` for querying resources across organizations


## [0.1.0] - 2025-02-11

Initial release of the HCP Terraform CLI with comprehensive API coverage.

### Added

- **Automated releases**: GoReleaser workflow for creating GitHub releases with binaries
  - Multi-platform binaries: Linux (amd64, arm64), macOS (Intel, Apple Silicon), Windows, FreeBSD
  - Automatic changelog generation from commits
  - SHA256 checksums for all binaries
  - Release process and checksums are published on GitHub Releases

- **TFE_ADDRESS support**: Backward compatibility with Terraform Enterprise
  - Environment variable precedence: HCPTF_ADDRESS > TFE_ADDRESS > default
  - Seamless migration for existing TFE users
  - See [docs/AUTH_GUIDE.md](docs/AUTH_GUIDE.md) for details

- **Comprehensive testing infrastructure**:
  - Added internal/client tests (92.9% coverage)
  - Added internal test guidance following HashiCorp go-tfe patterns
  - Added docs/TEST_COVERAGE.md with coverage report and roadmap
  - Internal packages average 79% coverage
  - Fixed broken test files and added test helpers

- **URL-style navigation**: Access resources using path-like syntax (e.g., `hcptf my-org my-workspace runs list`)
  - `hcptf <org>` - Show organization details
  - `hcptf <org> workspaces` - List workspaces
  - `hcptf <org> <workspace>` - Show workspace details
  - `hcptf <org> <workspace> runs` - List runs
  - `hcptf <org> <workspace> runs <run-id> <action>` - Perform run actions
  - See [docs/URL_NAVIGATION.md](docs/URL_NAVIGATION.md) for complete guide
  - Traditional command syntax remains fully supported

- Workspace commands: `workspace list|create|read|update|delete`
- Run commands: `run list|create|show|apply|discard|cancel`
- Organization commands: `organization list|create|show|update|delete`
- Variable commands: `variable list|create|update|delete`
- Team commands: `team list|create|show|delete|add-member|remove-member`
- Project commands: `project list|create|read|update|delete`
- State commands: `state list|read|outputs`
- Policy commands: `policy list|create|read|update|delete`
- Policy set commands: `policyset list|create|read|update|delete|add-policy|remove-policy`
- SSH key commands: `sshkey list|create|read|update|delete`
- Notification commands: `notification list|create|read|update|delete|verify`
- Variable set commands: `variableset list|create|read|update|delete|apply` and `variableset variable list|create|update|delete`
- Agent pool commands: `agentpool list|create|read|update|delete|token-create|token-list|token-delete`
- Run task commands: `runtask list|create|read|update|delete|attach|detach`
- OAuth client commands: `oauthclient list|create|read|update|delete`
- OAuth token commands: `oauthtoken list|read|update`
- Run trigger commands: `runtrigger list|create|read|delete`
- Plan commands: `plan read|logs`
- Apply commands: `apply read|logs`
- Configuration version commands: `configversion list|create|read|upload`
- Team access commands: `teamaccess list|create|read|update|delete`
- Project team access commands: `projectteamaccess list|create|read|update|delete`
- Registry module commands: `registrymodule list|create|read|delete|create-version|delete-version`
- Registry provider commands: `registryprovider list|create|read|delete`
- Registry provider version commands: `registryproviderversion create|read|delete`
- Registry provider platform commands: `registryproviderplatform create|read|delete`
- GPG key commands: `gpgkey list|create|read|update|delete`
- Stack commands: `stack list|create|read|update|delete`
- Stack configuration commands: `stackconfiguration list|create|read|update|delete`
- Stack deployment commands: `stackdeployment list|create|read`
- Stack state commands: `stackstate list|read`
- Audit trail commands: `audittrail list|read`
- Audit trail token commands: `audittrailtoken list|create|read|delete`
- Organization token commands: `organizationtoken list|create|read|delete`
- User token commands: `usertoken list|create|read|delete`
- Team token commands: `teamtoken list|create|read|delete`
- Organization membership commands: `organizationmembership list|create|read|delete`
- Organization member commands: `organizationmember read`
- Organization tag commands: `organizationtag list|delete`
- Reserved tag key commands: `reservedtagkey list|create|delete`
- Comment commands: `comment list|create|read`
- Policy check commands: `policycheck list|read|override`
- Policy evaluation commands: `policyevaluation list`
- Policy set outcome commands: `policysetoutcome list|read`
- Policy set parameter commands: `policysetparameter list|create|update|delete`
- AWS OIDC commands: `awsoidc create|read|update|delete`
- Azure OIDC commands: `azureoidc create|read|update|delete`
- GCP OIDC commands: `gcpoidc create|read|update|delete`
- Vault OIDC commands: `vaultoidc create|read|update|delete`
- Workspace resource commands: `workspaceresource list|read`
- Workspace tag commands: `workspacetag list|add|remove`
- Query run commands: `queryrun list`
- Query workspace commands: `queryworkspace list`
- Change request commands: `changerequest list|create|read|update`
- Assessment result commands: `assessmentresult list|read`
- HYOK configuration commands: `hyok list|create|read|update|delete`
- HYOK key version commands: `hyokkey create|read|delete`
- VCS event commands: `vcsevent list|read`
- Plan export commands: `planexport create|read|download`
- Agent monitoring commands: `agent list|read`
- Account commands: `account create|show|update`
- `login` and `logout` with token validation and Terraform CLI-compatible credential storage
- `version` command
- Multi-source authentication: `TFE_TOKEN`, `HCPTF_TOKEN`, `~/.hcptfrc`, Terraform CLI credentials
- Table and JSON output formats
- HCL-based configuration (`~/.hcptfrc`)

Total: 229 commands across 59 resource types.

[Unreleased]: https://github.com/thrashr888/hcptf-cli/compare/v0.5.0...HEAD
[0.5.0]: https://github.com/thrashr888/hcptf-cli/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/thrashr888/hcptf-cli/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/thrashr888/hcptf-cli/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/thrashr888/hcptf-cli/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/thrashr888/hcptf-cli/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/thrashr888/hcptf-cli/releases/tag/v0.1.0
