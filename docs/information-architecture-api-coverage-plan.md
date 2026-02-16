# Information Architecture + API Coverage Plan

Last updated: 2026-02-15

## Goal
Create a single, coherent CLI information architecture that:
- supports nested hierarchical commands as the canonical UX,
- preserves URL-style navigation as a first-class path,
- and reaches full practical API coverage (CRUD + action endpoints).

## Problem Summary
- Command registration, namespace help, router translation, and API coverage tracking are partially decoupled.
- Nested command support exists in some domains (`registry`, `stack`, `publicregistry`) but is inconsistent across the CLI.
- URL-style translation in `internal/router/router.go` relies on hardcoded command/resource lists that drift from registered commands.
- Coverage tracking currently uses filename heuristics and stale mappings, which hides real progress and creates false positives/negatives.

## Target Architecture
Use one source of truth for operations and command topology.

### Canonical model
- `Resource`: canonical noun path (`organization/member`, `policyset/parameter`, `stack/deployment`).
- `Operation`: `method`, `api_path`, `scope`, `verb`, `command_path`, `aliases`.
- `CommandPath`: tokenized hierarchical path (`policyset parameter list`).
- `URLRoutePattern`: optional URL-style route that resolves to the same operation.

### Outputs generated from the model
- Command registry (including parent namespaces at each depth).
- Namespace synopses/help lists.
- URL-style router known-command/resource sets and translation rules.
- Coverage matrix and progress report.

## Phased Execution Plan

### Phase 0: Baseline + Guardrails
Deliverables:
- Freeze current command key snapshot from `command/commands.go`.
- Add a CI-friendly check that detects router/command registry drift.

Acceptance:
- `go run . <namespace> -h` never misroutes to org/workspace URL context for known namespaces.

### Phase 1: Canonical Command Tree
Deliverables:
- Introduce an internal command tree abstraction (can start as static Go structures).
- Represent all current command paths (including nested paths already in use).
- Add canonical nested names for flattened namespaces, with aliases for backwards compatibility.

Examples:
- Canonical: `organization member read`
- Alias (compat): `organizationmember read`

Acceptance:
- Both canonical and alias forms resolve to the same command handler.
- Root and namespace help show canonical nested structure.

### Phase 2: Router Refactor (URL-style V2)
Deliverables:
- Replace hardcoded `knownCommands`/resource keyword lists with command-tree-derived resolution.
- Ensure help flags (`-h`, `--help`) always produce contextual help and never trigger network calls.
- Support URL-style for CRUD and action-style paths where unambiguous.

Acceptance:
- URL-style and traditional command-style are functionally equivalent for covered operations.
- Ambiguous URL paths fail with deterministic guidance, not accidental execution.

### Phase 3: API Operation Catalog
Deliverables:
- Build a catalog of API endpoints (method + path + operation semantic).
- Map each endpoint to canonical command path and verb (`list/read/create/update/delete` or explicit action).
- Track status: `DONE`, `TODO`, `IGNORE` with rationale.

Acceptance:
- Every endpoint in scope has a catalog entry and command mapping decision.

### Phase 4: Implementation Wave for Missing Coverage
Deliverables:
- Implement missing commands by priority:
1. High-use CRUD endpoints
2. Action endpoints with operational value (`override`, `safe-delete`, lifecycle actions)
3. Relationship endpoints (tag/member link/unlink semantics)
- Ensure each command is registered in canonical namespace and aliases where needed.

Acceptance:
- Operation catalog `TODO` count decreases each sprint.
- New commands include tests for flags, endpoint/method dispatch, and output format.

### Phase 5: Docs + UX Alignment
Deliverables:
- Update README examples to canonical nested syntax plus URL-style equivalents.
- Document alias policy and deprecation timeline for legacy flattened names.
- Ensure all namespace help text is generated or aligned from the same model.

Acceptance:
- README examples execute as documented.
- Help output and routing behavior are consistent.

## Namespace Normalization Strategy

### Keep as-is (already clear)
- `workspace`, `run`, `organization`, `project`, `team`, `stack`, `registry`, `publicregistry`

### Namespace normalization status

#### Done
- [x] `organizationmember` -> canonical `organization member`
- [x] `organizationmembership` -> canonical `organization membership`
- [x] `organizationtoken` -> canonical `organization token`
- [x] `policysetparameter` -> canonical `policyset parameter`
- [x] `policysetoutcome` -> canonical `policyset outcome`
- [x] `projectteamaccess` -> canonical `project teamaccess`
- [x] `teamaccess` -> canonical `team access`
- [x] `workspaceresource` -> canonical `workspace resource`
- [x] `workspacetag` -> canonical `workspace tag`

#### Pending
- [ ] None

## API Coverage Rules

### Verb mapping defaults
- `GET collection` -> `list`
- `GET item` -> `read`
- `POST` -> `create` unless endpoint is action-oriented
- `PATCH` -> `update`
- `DELETE` -> `delete`

### GET intent inference policy
- Allow optional omission of `list`/`read` only for non-mutating GET operations.
- Resolution order:
- explicit verb wins (`list`, `read`, `show`)
- if identity selector exists (`-id`, `-name`, canonical positional ID), infer `read`
- if only collection scope exists (`-org`, `-workspace`, parent scope), infer `list`
- if a resource has exactly one GET operation, infer that operation
- if ambiguous, fail with deterministic suggestion and do not issue API calls
- Never infer mutating or action verbs (`create`, `update`, `delete`, `apply`, `override`, etc.).

Concise examples:
- `hcptf workspace -org=my-org` → infer `list`
- `hcptf workspace -org=my-org -name=staging` → infer `read`
- `hcptf run -id=run-abc123` → infer `read`
- `hcptf run show -id=run-abc123` → explicit legacy alias for `read`

### Action endpoint policy
- Preserve explicit action names where meaningful:
- Examples: `override`, `apply`, `cancel`, `discard`, `safe-delete`, `restore`
- Do not force action endpoints into CRUD verbs.

### Relationship endpoint policy
- Prefer semantic verbs over transport verbs:
- Example: `workspace tag add/remove` over raw POST/DELETE relationships.

### Destructive operation safety policy
- All `delete` commands must require one of:
- explicit bypass flag (`-f`, `-force`, or `-y`)
- interactive confirmation prompt accepted with exact `yes`
- If neither condition is met, command must not perform deletion.
- This applies to all delete-style actions that remove resources or invalidate credentials/tokens.
- Help text for delete commands must clearly document:
- required identifier flags
- confirmation behavior
- bypass flags and their risk

Example:
`hcptf workspace delete -org=my-org -name=staging -y` (bypass prompt)

## Testing Strategy

### Unit tests
- Command parser tests for canonical + alias forms.
- Router translation tests for URL-style to canonical command path.
- Coverage script tests (or snapshot tests) against known command registry samples.
- Delete command tests must include guardrail coverage:
- requires identifier flags
- prompt/bypass behavior (`-f`/`-y`)
- cancellation path (no deletion performed)

### Integration-level checks
- `go run . -h` includes expected canonical namespaces.
- `go run . <known-namespace> -h` never routes into org/workspace context helpers.
- URL-style help cases:
- `<org> -h`
- `<org> <workspace> -h`
- `<org> <workspace> runs <run-id> <action> -h`

### Coverage gating
- Coverage report must consume command registry keys, not filename heuristics alone.
- CI should fail if router known commands diverge from registered namespace roots.

## Migration and Compatibility
- Keep legacy command aliases until users have migrated to canonical nested forms.
- Mark aliases in help output as compatibility paths.
- Remove aliases only after:
- 2 release cycles with warnings
- docs and examples fully migrated
- telemetry/user feedback confirms low dependency

## Immediate Next Milestones
1. Fix coverage tooling to be command-registry-aware (done in `scripts/api-coverage.sh` update).
2. Implement router/known-command synchronization.
3. Introduce canonical nested aliases for first high-confusion namespaces.
4. Start operation-catalog-driven implementation sprints.
