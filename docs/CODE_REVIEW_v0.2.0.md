# Code Quality Review: Changes Since v0.2.0

**Date:** 2026-02-15
**Scope:** 8 commits, 77 files changed, ~5,700 lines added, ~1,100 removed
**Build status:** `go test ./...` PASS, `go vet ./...` clean

## Overall Assessment

The changes add significant new command coverage and improve router architecture. The code follows consistent patterns and all tests pass. However, there are a few bugs that should be fixed before release, and the test coverage for new commands is thin.

---

## HIGH Severity

### 1. Missing `assessment` in short-form run-ID dispatch (BUG)

**`internal/router/router.go:108-148`**

The long-form path (`org workspace runs run-id assessment`) correctly maps to `assessmentresult list`, but the short-form path (`org workspace run-abc123 assessment`) has no `assessment` case. It falls through to produce `["run", "assessment", "-id=run-abc123"]`, which is not a valid command. This is a functional bug caused by the two dispatch blocks drifting out of sync.

**Fix:** Add an `"assessment"` case in the short-form block (around line 141):

```go
if action == "assessment" {
    return appendRemaining([]string{"assessmentresult", "list", "-org=" + org, "-workspace=" + workspace}, args, 4), nil
}
```

### 2. Unguarded type assertion can panic in `policycheck_read.go`

**`command/policycheck_read.go:148,164`**

```go
if !result.(bool) {
```

`result` is `interface{}` from a JSON map. If the field is `nil`, missing, or non-boolean, this panics at runtime. Should use the comma-ok idiom:

```go
resultBool, ok := result.(bool)
if !ok {
    // handle gracefully
}
```

### 3. `registrymodule_delete.go` `-provider` flag is dead code

**`command/registrymodule_delete.go:68`**

The `-provider` flag is accepted and shown in help text but never passed to the `Delete()` API call. The module is always fully deleted regardless of provider. Users could be misled into thinking they are deleting only a specific provider variant, potentially causing unintended data loss.

---

## MEDIUM Severity

### 4. No `url.PathEscape()` on user-supplied URL path segments

**10+ new command files (`githubapp_list.go`, `nocode_create.go`, `costestimate_read.go`, `githubapp_read.go`, `nocode_list.go`, `nocode_read.go`, `nocode_update.go`, `organizationtag_create.go`, `subscription_read.go`, `user_read.go`)**

Organization names, IDs, and other user inputs are interpolated directly into URL paths via `fmt.Sprintf` without escaping. While the server should reject malformed paths, the client should sanitize with `url.PathEscape()`.

### 5. Duplicated run-action dispatch -- already drifted

**`internal/router/router.go:108-148` and `166-204`**

The same run sub-resource mapping is duplicated between the short-form and long-form paths. They have already diverged (issue #1 above). Should be extracted into a shared helper:

```go
func (r *Router) translateRunAction(org, workspace, runID, action string, remaining []string) []string
```

### 6. Dead code branch in `inferImplicitGetVerb`

**`main.go:181-187`**

```go
remaining := tokens[matchLen:]
if len(remaining) > 0 {
    if isExplicitGetVerb(remaining[0]) {
        return args, nil
    }
    return args, nil  // <-- same result regardless of inner if
}
```

Both branches of the `isExplicitGetVerb` check return `args, nil` -- the conditional has no effect. Either simplify or implement the intended differentiation.

### 7. `nocode_read` and `nocode_list` are functionally identical

**`command/nocode_read.go` and `command/nocode_list.go`**

Both hit the same endpoint (`/api/v2/organizations/%s/no-code-provisioning`) with identical logic. `nocode_read` should presumably take an ID for a specific resource.

### 8. Test model has `"team read"` but registry has `"team show"`

**`internal/router/router_test.go:49`**

`testCommandPaths` uses `"team read"` but the real registry has `"team show"`. The guard test only checks root tokens, so this divergence is masked. Also missing `"team add-member"` and `"team remove-member"`.

### 9. `organizationtag_create.go` uses `%q` for JSON construction

**`command/organizationtag_create.go:48`**

```go
body := fmt.Sprintf(`{"data":{"type":"tags","attributes":{"name":%q}}}`, c.name)
```

Go's `%q` escaping differs from JSON escaping (e.g., `\x` escapes aren't valid JSON). Should use `json.Marshal()` or a proper struct.

### 10. README command table is stale

**`README.md`**

Missing entire command groups: `costestimate`, `featureset`, `githubapp`, `iprange`, `nocode`, `stabilitypolicy`, `subscription`, `user`. Also has wrong counts:

| Resource | README says | Actual |
|----------|-----------|--------|
| `oauthtoken` | 3 | 4 (missing `delete`) |
| `planexport` | 3 | 4 (missing `delete`) |
| `organization tag` | 2 | 3 (missing `create`) |
| `reservedtagkey` | 3 | 4 (missing `update`) |

### 11. Help text in 5 delete commands uses legacy flat command names

**`command/registrymodule_delete.go:92`, `command/registryprovider_delete.go:102`, `command/registryproviderplatform_delete.go:129`, `command/registryproviderversion_delete.go:111`, `command/policysetparameter_delete.go:85`**

These reference old names like `hcptf registrymodule delete` while the changelog says legacy aliases have been removed. Canonical names are now `hcptf registry module delete`, etc.

### 12. Inconsistent parse-error handling across commands

Some commands output raw API response body on error to `Ui.Output` (stdout), others to `Ui.Error` (stderr), and others suppress it entirely. Error data should consistently go to the error stream.

- `githubapp_list.go:56` uses `Ui.Output` (wrong)
- `githubapp_read.go:55`, `costestimate_read.go:54` use `Ui.Error` (correct)
- `featureset_list.go`, `nocode_*.go`, `subscription_*.go`, etc. suppress raw body entirely

### 13. `policycheck_read.go` bypasses `Meta.NewFormatter()` pattern

**`command/policycheck_read.go:51-53`**

```go
if c.Meta.OutputWriter == nil && c.Meta.ErrorWriter == nil {
    formatter = output.NewFormatterWithWriters(c.format, os.Stdout, os.Stderr)
}
```

`Meta.NewFormatter()` already handles this fallback. This special case breaks in tests where only one writer is set.

---

## LOW Severity

### 14. Credentials directory created with `0755` instead of `0700`

**`internal/config/config.go:233`**

The test helper correctly uses `0700`, but production uses `0755`. The file itself is `0600`, but a world-readable directory reveals the credential file's existence.

### 15. `ValidateToken` passes `nil` context

**`internal/config/config.go:219`**

`client.Users.ReadCurrent(nil)` has no timeout. If the API server hangs, this blocks forever. Should use `context.Background()` or a context with a deadline.

### 16. New HTTP client on every raw API call

**`command/api_request_helpers.go:28`**

`newHTTPClient()` is called per-request, defeating connection pooling. The `http.Client` (or at minimum the `http.Transport`) should be reused.

### 17. `pluralize("key")` produces `"keies"` not `"keys"`

**`internal/router/command_tree.go:177-179`**

The `y -> ies` rule doesn't check for a preceding vowel. English only applies this rule when the preceding letter is a consonant (`policy` -> `policies` is correct, but `key` -> `keys`, not `keies`). Benign currently but a latent bug.

### 18. Thin test coverage for new httptest-based commands

11 of 14 new httptest-based test files lack any API error tests. Several files have only a single happy-path test:

| File | Tests | Error path | Grade |
|------|-------|------------|-------|
| `stabilitypolicy_read_test.go` | 1 | No | D |
| `subscription_list_test.go` | 1 | No | D |
| `iprange_list_test.go` | 1 | No | D |
| `costestimate_read_test.go` | 2 | No | C |
| `githubapp_list_test.go` | 2 | No | C |
| `githubapp_read_test.go` | 2 | No | C |
| `nocode_list_test.go` | 2 | No | C |
| `nocode_read_test.go` | 2 | No | C |
| `whoami_test.go` | 6 | Yes | A |
| `reservedtagkey_update_test.go` | 5 | Yes | A- |

`featureset_list_test.go:11-18` (`TestFeatureSetListCommand_RequiresNoArgs`) is a no-op test that only checks a struct literal is not nil.

### 19. No `-output` format validation

No command validates that `-output` is `"table"` or `"json"`. Invalid values like `-output=yaml` silently fall through.

### 20. Missing `Accept` and `User-Agent` headers in raw API requests

**`command/api_request_helpers.go`**

Sets `Content-Type: application/vnd.api+json` but not `Accept`. No `User-Agent` header is set.

### 21. `hasHelpFlag` duplicated in `router.go` and `main.go`

**`internal/router/router.go:300` and `main.go:272`**

Two identical implementations. Should be shared.

---

## Positives

- **Security:** No credential leaks found. Token is never logged or included in error messages. Credential file uses `0600` permissions. Delete commands all have confirmation prompts with test coverage.
- **Architecture:** The router/command-tree generation from the command model is well-designed and self-validating via the roots-match test.
- **Consistency:** New commands follow the established `Meta` pattern for formatting and CLI interaction.
- **Safety:** All 8 reviewed delete commands implement the 3-layer safety mechanism (`-force`, `-y`, interactive prompt) with tests covering each path.
- **No race conditions, no concurrency issues, no command injection vectors.**
- **No sensitive information in documentation** -- all examples use placeholder tokens.

---

## Recommended Fix Priority

| Priority | Issues | Effort |
|----------|--------|--------|
| **Before release** | #1 (assessment dispatch bug), #2 (panic), #3 (dead provider flag) | Small |
| **Soon after** | #4 (path escaping), #5 (extract run-action helper), #9 (JSON construction), #10 (README) | Medium |
| **When convenient** | #6-8, #11-21 | Varies |
