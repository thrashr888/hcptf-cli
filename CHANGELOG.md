# Changelog

All notable changes to the HCP Terraform CLI (`hcptf`).

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/thrashr888/hcptf-cli/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/thrashr888/hcptf-cli/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/thrashr888/hcptf-cli/releases/tag/v0.1.0
