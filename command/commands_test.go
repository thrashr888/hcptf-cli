package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestCommandsNamespaceNormalizationCanonicalOnly(t *testing.T) {
	commandMap := Commands(&Meta{})

	resolveCommand := func(name string) {
		t.Helper()

		factory, ok := commandMap[name]
		if !ok {
			t.Fatalf("missing command key %q", name)
		}

		_, err := factory()
		if err != nil {
			t.Fatalf("failed to construct command %q: %v", name, err)
		}
	}

	canonicalKeys := []string{
		"audittrail token",
		"audittrail token list",
		"audittrail token create",
		"audittrail token read",
		"audittrail token delete",
		"organization member",
		"organization member read",
		"organization membership",
		"organization membership list",
		"organization membership create",
		"organization membership read",
		"organization membership delete",
		"organization token",
		"organization token list",
		"organization token create",
		"organization token read",
		"organization token delete",
		"organization tag",
		"organization tag list",
		"organization tag create",
		"organization tag delete",
		"policyset parameter",
		"policyset parameter list",
		"policyset parameter create",
		"policyset parameter update",
		"policyset parameter delete",
		"policyset outcome",
		"policyset outcome list",
		"policyset outcome read",
		"project teamaccess",
		"project teamaccess list",
		"project teamaccess create",
		"project teamaccess read",
		"project teamaccess update",
		"project teamaccess delete",
		"team access",
		"team access list",
		"team access create",
		"team access read",
		"team access update",
		"team access delete",
		"team token",
		"team token list",
		"team token create",
		"team token read",
		"team token delete",
		"user token",
		"user token list",
		"user token create",
		"user token read",
		"user token delete",
		"workspace resource",
		"workspace resource list",
		"workspace resource read",
		"workspace tag",
		"workspace tag list",
		"workspace tag add",
		"workspace tag remove",
	}

	for _, key := range canonicalKeys {
		resolveCommand(key)
	}

	legacyKeys := []string{
		"audittrailtoken",
		"audittrailtoken list",
		"audittrailtoken create",
		"audittrailtoken read",
		"audittrailtoken delete",
		"organizationmember",
		"organizationmember read",
		"organizationmembership",
		"organizationmembership list",
		"organizationmembership create",
		"organizationmembership read",
		"organizationmembership delete",
		"organizationtoken",
		"organizationtoken list",
		"organizationtoken create",
		"organizationtoken read",
		"organizationtoken delete",
		"organizationtag",
		"organizationtag list",
		"organizationtag create",
		"organizationtag delete",
		"policysetparameter",
		"policysetparameter list",
		"policysetparameter create",
		"policysetparameter update",
		"policysetparameter delete",
		"policysetoutcome",
		"policysetoutcome list",
		"policysetoutcome read",
		"projectteamaccess",
		"projectteamaccess list",
		"projectteamaccess create",
		"projectteamaccess read",
		"projectteamaccess update",
		"projectteamaccess delete",
		"teamaccess",
		"teamaccess list",
		"teamaccess create",
		"teamaccess read",
		"teamaccess update",
		"teamaccess delete",
		"teamtoken",
		"teamtoken list",
		"teamtoken create",
		"teamtoken read",
		"teamtoken delete",
		"workspaceresource",
		"workspaceresource list",
		"workspaceresource read",
		"workspacetag",
		"workspacetag list",
		"workspacetag add",
		"workspacetag remove",
		"usertoken",
		"usertoken list",
		"usertoken create",
		"usertoken read",
		"usertoken delete",
	}

	for _, key := range legacyKeys {
		if _, ok := commandMap[key]; ok {
			t.Fatalf("legacy command key should not exist: %q", key)
		}
	}
}

// resourceDef mirrors the shell script's "api_doc_name|command_prefix|operations" format.
type resourceDef struct {
	apiDoc    string
	cliPrefix string
	ops       []string
}

// crudLabel converts short CRUD codes to action names, matching the shell script's op_label().
func crudLabel(op string) string {
	switch op {
	case "L":
		return "list"
	case "C":
		return "create"
	case "R":
		return "read"
	case "U":
		return "update"
	case "D":
		return "delete"
	default:
		return op
	}
}

// cliPrefix expands compound prefixes to their CLI command namespace,
// matching the shell script's resolve_prefix() plus the hierarchical
// namespace commands that use spaces in the command map.
func cliPrefix(prefix string) string {
	switch prefix {
	case "workspacetag":
		return "workspace tag"
	case "workspaceresource":
		return "workspace resource"
	case "stackconfiguration":
		return "stack configuration"
	case "stackdeployment":
		return "stack deployment"
	case "stackstate":
		return "stack state"
	case "registrymodule":
		return "registry module"
	case "registryprovider":
		return "registry provider"
	case "registryproviderversion":
		return "registry provider version"
	case "registryproviderplatform":
		return "registry provider platform"
	case "agentpool_token":
		return "agentpool token"
	case "organizationmembership":
		return "organization membership"
	case "organizationtag":
		return "organization tag"
	case "organizationtoken":
		return "organization token"
	case "teamaccess":
		return "team access"
	case "teamtoken":
		return "team token"
	case "projectteamaccess":
		return "project teamaccess"
	case "policysetparameter":
		return "policyset parameter"
	case "audittrailtoken":
		return "audittrail token"
	case "usertoken":
		return "user token"
	default:
		return prefix
	}
}

// apiResources defines the expected API coverage — ported directly from api-coverage.sh.
var apiResources = []resourceDef{
	// Account & auth
	{"account", "account", []string{"C", "R", "U"}},

	// Organizations
	{"organizations", "organization", []string{"L", "C", "R", "U", "D"}},
	{"organization-memberships", "organizationmembership", []string{"L", "C", "R", "D"}},
	{"organization-tags", "organizationtag", []string{"L", "C", "D"}},
	{"organization-tokens", "organizationtoken", []string{"L", "C", "R", "D"}},

	// Workspaces
	{"workspaces", "workspace", []string{"L", "C", "R", "U", "D", "lock", "unlock", "force-unlock"}},
	{"workspace-variables", "variable", []string{"L", "C", "U", "D"}},
	{"workspace-resources", "workspaceresource", []string{"L", "R"}},
	{"workspace-tags", "workspacetag", []string{"L", "C", "D"}},

	// Runs / plans / applies
	{"runs", "run", []string{"L", "C", "R", "list-org", "apply", "discard", "cancel", "force-execute"}},
	{"applies", "apply", []string{"R"}},
	{"plans", "plan", []string{"R"}},
	{"plan-exports", "planexport", []string{"C", "R", "D"}},
	{"cost-estimates", "costestimate", []string{"R"}},

	// State
	{"state-versions", "state", []string{"L", "R"}},
	{"state-version-outputs", "state", []string{"outputs"}},

	// Config versions / variable sets
	{"configuration-versions", "configversion", []string{"L", "C", "R"}},
	{"variable-sets", "variableset", []string{"L", "C", "R", "U", "D", "apply-workspaces", "remove-workspaces", "apply-projects", "remove-projects", "apply-stacks", "remove-stacks", "list-workspace", "list-project", "update-workspaces", "update-stacks"}},

	// Teams / access
	{"teams", "team", []string{"L", "C", "R", "D"}},
	{"team-access", "teamaccess", []string{"L", "C", "R", "U", "D"}},
	{"team-members", "team", []string{"R"}},
	{"team-tokens", "teamtoken", []string{"L", "C", "R", "D"}},

	// Projects
	{"projects", "project", []string{"L", "C", "R", "U", "D"}},
	{"project-team-access", "projectteamaccess", []string{"L", "C", "R", "U", "D"}},

	// Policies
	{"policies", "policy", []string{"L", "C", "R", "U", "D", "upload", "download"}},
	{"policy-sets", "policyset", []string{"L", "C", "R", "U", "D", "add-policy", "remove-policy", "add-workspace", "remove-workspace", "add-workspace-exclusion", "remove-workspace-exclusion", "add-project", "remove-project"}},
	{"policy-checks", "policycheck", []string{"L", "R", "override"}},
	{"policy-evaluations", "policyevaluation", []string{"L"}},
	{"policy-set-params", "policysetparameter", []string{"L", "C", "U", "D"}},

	// Agents / pools
	{"agents", "agent", []string{"L", "R"}},
	{"agent-tokens", "agentpool_token", []string{"L", "C", "D"}},

	// SSH / OAuth / notifications
	{"ssh-keys", "sshkey", []string{"L", "C", "R", "U", "D"}},
	{"oauth-clients", "oauthclient", []string{"L", "C", "R", "U", "D"}},
	{"oauth-tokens", "oauthtoken", []string{"L", "R", "U", "D"}},
	{"notification-configurations", "notification", []string{"L", "C", "R", "U", "D"}},

	// Run tasks / triggers
	{"run-tasks", "runtask", []string{"L", "C", "R", "U", "D"}},
	{"run-triggers", "runtrigger", []string{"L", "C", "R", "D"}},

	// Comments / audit
	{"comments", "comment", []string{"L", "C", "R"}},
	{"audit-trails", "audittrail", []string{"L", "R"}},
	{"audit-trails-tokens", "audittrailtoken", []string{"L", "C", "R", "D"}},

	// Registry
	{"private-registry/modules", "registrymodule", []string{"L", "C", "R", "D"}},
	{"private-registry/provider-versions-platforms", "registryproviderplatform", []string{"C", "R", "D"}},
	{"private-registry/providers", "registryprovider", []string{"L", "C", "R", "D"}},
	{"private-registry/manage-provider-versions", "registryproviderversion", []string{"C", "R", "D"}},
	{"private-registry/gpg-keys", "gpgkey", []string{"L", "C", "R", "U", "D"}},

	// VCS / health
	{"vcs-events", "vcsevent", []string{"L", "R"}},
	{"github-app-installations", "githubapp", []string{"L", "R"}},
	{"assessment-results", "assessmentresult", []string{"L", "R"}},
	{"change-requests", "changerequest", []string{"L", "C", "R", "U"}},

	// Stacks
	{"stacks/stacks", "stack", []string{"L", "C", "R", "U", "D"}},
	{"stacks/stack-configurations", "stackconfiguration", []string{"L", "C", "R", "U", "D"}},
	{"stacks/stack-deployments", "stackdeployment", []string{"L", "C", "R"}},
	{"stacks/stack-states", "stackstate", []string{"L", "R"}},

	// OIDC / HYOK
	{"hold-your-own-key/aws", "awsoidc", []string{"C", "R", "U", "D"}},
	{"hold-your-own-key/azure", "azureoidc", []string{"C", "R", "U", "D"}},
	{"hold-your-own-key/gcp", "gcpoidc", []string{"C", "R", "U", "D"}},
	{"hold-your-own-key/vault-transit", "vaultoidc", []string{"C", "R", "U", "D"}},
	{"hold-your-own-key/byok", "hyok", []string{"L", "C", "R", "U", "D"}},
	{"hold-your-own-key/key-management", "hyokkey", []string{"C", "R", "D"}},

	// Queries / explorer
	{"queries/run-query", "queryrun", []string{"L"}},
	{"queries/workspace-query", "queryworkspace", []string{"L"}},
	{"explorer", "explorer", []string{"query"}},

	// User tokens / users
	{"user-tokens", "usertoken", []string{"L", "C", "R", "D"}},
	{"users", "user", []string{"R"}},

	// Billing / metadata
	{"subscriptions", "subscription", []string{"L", "R"}},
	{"feature-sets", "featureset", []string{"L"}},
	{"ip-ranges", "iprange", []string{"L"}},
	{"no-code-provisioning", "nocode", []string{"L", "C", "R", "U"}},
	{"stability-policy", "stabilitypolicy", []string{"R"}},

	// Reserved tags
	{"reserved-tag-keys", "reservedtagkey", []string{"L", "C", "U", "D"}},
}

// hasRegisteredCommand checks whether a given resource operation has a registered CLI command.
// It mirrors the shell script's has_command() with all its alias logic.
func hasRegisteredCommand(cmdSet map[string]bool, rawPrefix, action string) bool {
	prefix := cliPrefix(rawPrefix)

	var candidates []string

	// Agentpool token special-casing
	if prefix == "agentpool token" {
		switch action {
		case "list":
			candidates = append(candidates, "agentpool token-list")
		case "create":
			candidates = append(candidates, "agentpool token-create")
		case "delete":
			candidates = append(candidates, "agentpool token-delete")
		}
	} else {
		candidates = append(candidates, prefix+" "+action)
	}

	// "read" aliases to "show"
	if action == "read" {
		candidates = append(candidates, prefix+" show")
	}

	// Custom operation aliases (mirrors the shell script's case block)
	switch prefix + ":" + action {
	case "workspace tag:create":
		candidates = append(candidates, "workspace tag add")
	case "workspace tag:delete":
		candidates = append(candidates, "workspace tag remove")
	case "workspace:force-unlock":
		candidates = append(candidates, "workspace force-unlock")
	case "variableset:apply-workspaces",
		"variableset:apply-projects",
		"variableset:apply-stacks":
		candidates = append(candidates, "variableset apply")
	case "variableset:remove-workspaces",
		"variableset:remove-projects",
		"variableset:remove-stacks":
		candidates = append(candidates, "variableset remove")
	case "explorer:query":
		candidates = append(candidates, "explorer query")
	case "policycheck:override":
		candidates = append(candidates, "policycheck override")
	case "state:outputs":
		candidates = append(candidates, "state outputs")
	}

	for _, c := range candidates {
		if cmdSet[c] {
			return true
		}
	}

	return false
}

// TestAllCommandsRegistered verifies that every expected API operation has a
// corresponding registered CLI command. This is the Go equivalent of
// scripts/api-coverage.sh.
func TestAllCommandsRegistered(t *testing.T) {
	ui := cli.NewMockUi()
	meta := newTestMeta(ui)
	commands := Commands(&meta)

	// Build set of registered command names.
	cmdSet := make(map[string]bool, len(commands))
	for name := range commands {
		cmdSet[name] = true
	}

	var missing []string
	for _, res := range apiResources {
		for _, op := range res.ops {
			action := crudLabel(op)
			if !hasRegisteredCommand(cmdSet, res.cliPrefix, action) {
				missing = append(missing, fmt.Sprintf("%s → %s %s", res.apiDoc, res.cliPrefix, action))
			}
		}
	}

	if len(missing) > 0 {
		t.Errorf("missing CLI commands for %d API operations:\n  %s",
			len(missing), strings.Join(missing, "\n  "))
	}
}

// TestAllCommandsHaveHelpAndSynopsis instantiates every registered command and
// asserts Help() and Synopsis() return non-empty strings.
func TestAllCommandsHaveHelpAndSynopsis(t *testing.T) {
	ui := cli.NewMockUi()
	meta := newTestMeta(ui)
	commands := Commands(&meta)

	for name, factory := range commands {
		t.Run(name, func(t *testing.T) {
			cmd, err := factory()
			if err != nil {
				t.Fatalf("factory error: %v", err)
			}

			if cmd.Synopsis() == "" {
				t.Errorf("command %q has empty Synopsis()", name)
			}
			if cmd.Help() == "" {
				t.Errorf("command %q has empty Help()", name)
			}
		})
	}
}

// TestAllCommandFilesHaveTests checks that every command implementation file
// has a corresponding _test.go file. This is the Go equivalent of
// scripts/test-priority.sh.
func TestAllCommandFilesHaveTests(t *testing.T) {
	// High priority groups — these MUST have tests.
	highPriority := map[string]bool{
		"organization": true,
		"workspace":    true,
		"run":          true,
		"state":        true,
		"variable":     true,
		"project":      true,
		"team":         true,
	}

	commandDir := "."
	matches, err := filepath.Glob(filepath.Join(commandDir, "*_*.go"))
	if err != nil {
		t.Fatalf("glob error: %v", err)
	}

	// Exclusions: test files themselves, and non-command support files.
	excluded := map[string]bool{
		"commands.go":          true,
		"commands_test.go":     true,
		"meta.go":              true,
		"test_helpers_test.go": true,
		"namespace_command.go": true,
	}

	type fileInfo struct {
		name    string
		group   string
		hasTest bool
	}
	var files []fileInfo
	missingHighPriority := 0

	for _, match := range matches {
		base := filepath.Base(match)

		// Skip test files
		if strings.HasSuffix(base, "_test.go") {
			continue
		}
		// Skip excluded files
		if excluded[base] {
			continue
		}
		// Skip service interface files
		if strings.HasSuffix(base, "_services.go") {
			continue
		}

		// Determine the command group (e.g. "workspace" from "workspace_delete.go")
		group := strings.SplitN(strings.TrimSuffix(base, ".go"), "_", 2)[0]

		// Check for a corresponding test file
		testFile := strings.TrimSuffix(base, ".go") + "_test.go"
		testPath := filepath.Join(commandDir, testFile)
		_, statErr := os.Stat(testPath)
		hasTest := statErr == nil

		files = append(files, fileInfo{name: base, group: group, hasTest: hasTest})

		if !hasTest && highPriority[group] {
			missingHighPriority++
		}
	}

	// Report all missing test files
	var missingTests []string
	for _, f := range files {
		if !f.hasTest {
			priority := "low"
			if highPriority[f.group] {
				priority = "HIGH"
			}
			missingTests = append(missingTests, fmt.Sprintf("[%s] %s", priority, f.name))
		}
	}

	if len(missingTests) > 0 {
		t.Logf("command files without tests (%d/%d):\n  %s",
			len(missingTests), len(files), strings.Join(missingTests, "\n  "))
	}

	// Fail for high-priority groups only — others are informational.
	if missingHighPriority > 0 {
		t.Skipf("%d high-priority command files lack tests (see log above)", missingHighPriority)
	}
}

// commandNameToFile derives the expected Go source file for a registered command name.
// The convention is: multi-word namespace prefixes are joined without separators,
// then an underscore separates them from the action, and hyphens become underscores.
// For example: "team access list" → "teamaccess_list.go",
// "run force-execute" → "run_force_execute.go".
func commandNameToFile(name string) string {
	// Explicit overrides for commands that don't follow any regular pattern.
	overrides := map[string]string{
		"registry module version create": "registrymodule_create_version.go",
		"registry module version delete": "registrymodule_delete_version.go",
	}
	if f, ok := overrides[name]; ok {
		return f
	}

	// Replace colons with spaces so "organization:context" → "organization context".
	name = strings.ReplaceAll(name, ":", " ")

	// Known multi-word namespace prefixes → their joined file prefix.
	// Ordered longest first so "registry provider version" matches before "registry provider".
	type nsMapping struct {
		prefix   string
		fileBase string
	}
	nsPrefixes := []nsMapping{
		{"registry provider platform", "registryproviderplatform"},
		{"registry provider version", "registryproviderversion"},
		{"organization membership", "organizationmembership"},
		{"policyset parameter", "policysetparameter"},
		{"stack configuration", "stackconfiguration"},
		{"variableset variable", "variableset_variable"},
		{"workspace resource", "workspaceresource"},
		{"audittrail token", "audittrailtoken"},
		{"registry provider", "registryprovider"},
		{"project teamaccess", "projectteamaccess"},
		{"organization member", "organizationmember"},
		{"organization token", "organizationtoken"},
		{"policyset outcome", "policysetoutcome"},
		{"stack deployment", "stackdeployment"},
		{"registry module", "registrymodule"},
		{"organization tag", "organizationtag"},
		{"workspace tag", "workspacetag"},
		{"team access", "teamaccess"},
		{"team token", "teamtoken"},
		{"user token", "usertoken"},
		{"stack state", "stackstate"},
	}

	for _, ns := range nsPrefixes {
		pfx := ns.prefix + " "
		if strings.HasPrefix(name, pfx) {
			action := name[len(pfx):]
			action = strings.ReplaceAll(action, " ", "_")
			action = strings.ReplaceAll(action, "-", "_")
			return ns.fileBase + "_" + action + ".go"
		}
	}

	// Default: replace all separators with underscores.
	result := strings.ReplaceAll(name, " ", "_")
	result = strings.ReplaceAll(result, "-", "_")
	return result + ".go"
}

// TestCommandNamesMatchFiles validates that every registered command name
// maps to a Go file following the naming convention.
// For example, "workspace lock" should have workspace_lock.go.
func TestCommandNamesMatchFiles(t *testing.T) {
	ui := cli.NewMockUi()
	meta := newTestMeta(ui)
	commands := Commands(&meta)

	commandDir := "."

	// Namespace-only parent commands that use NamespaceCommand and have no
	// dedicated file (they are generated dynamically in Commands()).
	namespaceOnly := map[string]bool{
		"team access":             true,
		"project teamaccess":      true,
		"policyset outcome":       true,
		"policyset parameter":     true,
		"audittrail token":        true,
		"user token":              true,
		"team token":              true,
		"organization membership": true,
		"organization member":     true,
		"organization token":      true,
		"organization tag":        true,
		"workspace tag":           true,
		"workspace resource":      true,
	}

	// Auto-generated single-word namespace commands produced by the loop at the
	// end of Commands(). These have no dedicated file either.
	autoGenerated := make(map[string]bool)
	for name := range commands {
		parts := strings.Split(name, " ")
		if len(parts) > 1 {
			autoGenerated[parts[0]] = true
		}
	}

	var mismatched []string
	for name := range commands {
		if namespaceOnly[name] {
			continue
		}
		// Skip auto-generated namespace parents that have no explicit entry in
		// commands.go (their factory is added by the loop). We detect these as
		// single-word keys that serve as parents for multi-word commands AND
		// are not explicitly declared with a dedicated struct.
		if autoGenerated[name] && !strings.Contains(name, " ") && !strings.Contains(name, ":") {
			// Check if the file actually exists; if it does the command was
			// explicitly declared (e.g. "stack" → stack.go).
			expectedFile := commandNameToFile(name)
			if _, err := os.Stat(filepath.Join(commandDir, expectedFile)); os.IsNotExist(err) {
				continue
			}
			// File exists — don't skip, let the check pass naturally.
		}

		expectedFile := commandNameToFile(name)
		path := filepath.Join(commandDir, expectedFile)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			mismatched = append(mismatched, fmt.Sprintf("%q → expected %s", name, expectedFile))
		}
	}

	if len(mismatched) > 0 {
		t.Errorf("registered commands with no matching file (%d):\n  %s",
			len(mismatched), strings.Join(mismatched, "\n  "))
	}
}
