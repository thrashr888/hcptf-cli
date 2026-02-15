# API Coverage Tracking

Last updated: 2026-02-15

This document tracks Terraform Cloud API coverage gaps against CLI command files in `/command/*_*.go`.

## Scope
- Source of truth for API docs: `content/terraform-docs-common/docs/cloud-docs/api-docs/*`
- Coverage target: CRUD + action style endpoints that should have corresponding CLI commands.
- Status meaning:
  - `DONE` = command file exists for that resource/operation pattern
  - `TODO` = missing command(s) to add
  - `IGNORE` = intentionally out of scope for now

## Quick conventions used
- `resource` maps to doc file name in API docs (e.g. `workspaces.mdx`).
- Commands are represented by existing filenames (e.g. `workspacetag_add`, `oauthclient_create`).

## High-priority missing resources (not currently represented)

- [ ] `cost-estimates`
  - API file: `cost-estimates.mdx`
  - Missing: `read` (GET /cost-estimates/:id)
- [ ] `feature-sets`
  - API file: `feature-sets.mdx`
  - Missing: `list` (GET /feature-sets)
- [ ] `github-app-installations`
  - API file: `github-app-installations.mdx`
  - Missing: `list` (+ detail/read variants around installations)
- [ ] `invoices`
  - API file: `invoices.mdx`
  - Missing: `list` (GET /organizations/:organization_name/invoices)
- [ ] `ip-ranges`
  - API file: `ip-ranges.mdx`
  - Missing: `list` (GET /meta/ip-ranges)
- [ ] `no-code-provisioning`
  - API file: `no-code-provisioning.mdx`
  - Missing: `create`, `read`, `update`
- [ ] `subscriptions`
  - API file: `subscriptions.mdx`
  - Missing: `list`, `read`
- [ ] `users`
  - API file: `users.mdx`
  - Missing: `read` (GET /users/:user_id)
- [ ] `users` and related auth-token views were intentionally added as TODO per your request.

## Partially covered resources (CLI exists, but missing op)

- [ ] `account` (`account.mdx`)
  - Missing: `read(list)` endpoint style for `/account/details`, and action endpoints on account updates/password
- [ ] `agent-tokens` (`agent-tokens.mdx`)
  - Missing: `create`, `delete`
- [ ] `agents` (`agents.mdx`)
  - Missing: `create`, `delete`, `update`
- [ ] `audit-trails-tokens` (`audit-trails-tokens.mdx`)
  - Missing: `action` for token revoke/rotation
- [ ] `change-requests` (`change-requests.mdx`)
  - Missing: `create` mapping to explorer bulk-actions action endpoint(s)
- [ ] `explorer` (`explorer.mdx`)
  - Missing: `create`, `list`, `action` (current CLI has only `explorer_query`)
- [ ] `oauth-tokens` (`oauth-tokens.mdx`)
  - Missing: `delete`
- [ ] `organization-tags` (`organization-tags.mdx`)
  - Missing: `create`, action-style attach/update paths (POST/DELETE relationship ops)
- [ ] `organization-tokens` (`organization-tokens.mdx`)
  - Missing: action endpoints for organization auth token management
- [ ] `organizations` (`organizations.mdx`)
  - Missing: extra action endpoint (data-retention-policy relationship operation)
- [ ] `plan-exports` (`plan-exports.mdx`)
  - Missing: `delete`
- [ ] `plans` (`plans.mdx`)
  - Missing: `list`, action-like plan output fetch endpoints
- [ ] `policy-checks` (`policy-checks.mdx`)
  - Missing: `create`/action endpoint (`/policy-checks/:id/actions/override`)
- [ ] `policy-evaluations` (`policy-evaluations.mdx`)
  - Missing: `read` for policy set outcome views
- [ ] `policy-set-params` (`policy-set-params.mdx`)
  - Missing: `read` (GET `/policy-sets/:policy_set_id/parameters`)
- [ ] `private-registry/gpg-keys` (`private-registry/gpg-keys.mdx`)
  - Missing: `action` endpoint style for registry key mutation semantics
- [ ] `private-registry/manage-module-versions` (`private-registry/manage-module-versions.mdx`)
  - Missing: `action` patch/update semantics
- [ ] `private-registry/modules` (`private-registry/modules.mdx`)
  - Missing: `action`/non-deprecated create/read/write flows
- [ ] `private-registry/test-configuration` (`private-registry/test-configuration.mdx`)
  - Missing: `read`, `update`
- [ ] `private-registry/tests` (`private-registry/tests.mdx`)
  - Missing: `create`, `list`, `action` (test run lifecycle)
- [ ] `project-team-access` and `projects` cross-ops (`projects.mdx`)
  - Missing: action route(s) for project tag/relationships operations
- [ ] `queries/index` (`queries/index.mdx`)
  - Missing: `create`, `read`
- [ ] `reserved-tag-keys` (`reserved-tag-keys.mdx`)
  - Missing: `update` on `/reserved-tags/:reserved_tag_key_id`
- [ ] `state-versions` (`state-versions.mdx`)
  - Missing: `create`, action routes (`soft_delete_backing_data`, `restore_backing_data`, `permanently_delete_backing_data`)
- [ ] `team-members` (`team-members.mdx`)
  - Missing: action-style relationship operations (POST/DELETE relationships)
- [ ] `workspaces` (`workspaces.mdx`)
  - Missing: action for safe-delete and relationship actions (`POST /workspaces/:id/actions/safe-delete`, tag/ssh-key relationships)

## Documentation-only/ambiguous entries (review before coding)
These files are mixed or index-style and should be validated manually before creating commands:

- [ ] `changelog.mdx` (aggregates many resource snippets)
- [ ] `index.mdx` (global docs index + auth/session endpoints)
- [ ] `_template.mdx` (non-API template doc)

## Out of scope (default Ignore)
Mark these as ignore if they should stay out of CLI work:

- [ ] `_template.mdx`
- [ ] `changelog.mdx`
- [ ] `index.mdx`
- [ ] Any deprecation-warning-only module endpoints under `private-registry/modules.mdx`

## Progress board
- Use this for sprint tracking:
  - `TODO` → implement
  - `DONE` → implemented and verified
  - `IGNORE` → intentionally out of scope

Suggested columns in an external tracker:
- `Resource`
- `Missing operations`
- `Command(s) to add`
- `Owner`
- `Status`
- `Notes`
