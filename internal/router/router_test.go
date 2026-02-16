package router

import (
	"reflect"
	"testing"
)

var testCommandPaths = []string{
	"version",
	"login",
	"logout",
	"whoami",
	"workspace list",
	"workspace read",
	"workspace create",
	"workspace update",
	"workspace delete",
	"run list",
	"run show",
	"run apply",
	"organization show",
	"organization list",
	"project list",
	"team list",
	"policy list",
	"policyset list",
	"variable list",
	"state list",
	"state outputs",
	"configversion list",
	"workspaceresource list",
	"workspacetag list",
	"comment list",
	"policycheck list",
	"assessmentresult list",
	"changerequest list",
	"organization:context",
	"workspace:context",
}

func newTestRouter() *Router {
	return NewRouter(nil, testCommandPaths)
}

func TestTranslateArgs(t *testing.T) {
	r := newTestRouter()

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty args",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "flag arg (passthrough)",
			input:    []string{"-h"},
			expected: []string{"-h"},
		},
		{
			name:     "known command (passthrough)",
			input:    []string{"workspace", "list", "-org=myorg"},
			expected: []string{"workspace", "list", "-org=myorg"},
		},
		{
			name:     "whoami command (passthrough)",
			input:    []string{"whoami", "-output=json"},
			expected: []string{"whoami", "-output=json"},
		},
		{
			name:     "org only",
			input:    []string{"myorg"},
			expected: []string{"organization", "show", "-name=myorg"},
		},
		{
			name:     "org workspaces",
			input:    []string{"myorg", "workspaces"},
			expected: []string{"workspace", "list", "-org=myorg"},
		},
		{
			name:     "org projects",
			input:    []string{"myorg", "projects"},
			expected: []string{"project", "list", "-org=myorg"},
		},
		{
			name:     "org teams",
			input:    []string{"myorg", "teams"},
			expected: []string{"team", "list", "-org=myorg"},
		},
		{
			name:     "org teams help",
			input:    []string{"myorg", "teams", "-h"},
			expected: []string{"team", "-h"},
		},
		{
			name:     "org teams create",
			input:    []string{"myorg", "teams", "create"},
			expected: []string{"team", "create", "-org=myorg"},
		},
		{
			name:     "org teams create help",
			input:    []string{"myorg", "teams", "create", "-h"},
			expected: []string{"team", "create", "-org=myorg", "-h"},
		},
		{
			name:     "org policies",
			input:    []string{"myorg", "policies"},
			expected: []string{"policy", "list", "-org=myorg"},
		},
		{
			name:     "org policysets",
			input:    []string{"myorg", "policysets"},
			expected: []string{"policyset", "list", "-org=myorg"},
		},
		{
			name:     "org workspace",
			input:    []string{"myorg", "myworkspace"},
			expected: []string{"workspace", "read", "-org=myorg", "-name=myworkspace"},
		},
		{
			name:     "org workspace runs",
			input:    []string{"myorg", "myworkspace", "runs"},
			expected: []string{"run", "list", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace runs list",
			input:    []string{"myorg", "myworkspace", "runs", "list"},
			expected: []string{"run", "list", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace runs show",
			input:    []string{"myorg", "myworkspace", "runs", "run-123"},
			expected: []string{"run", "show", "-id=run-123"},
		},
		{
			name:     "org workspace runs apply",
			input:    []string{"myorg", "myworkspace", "runs", "run-123", "apply"},
			expected: []string{"run", "apply", "-id=run-123"},
		},
		{
			name:     "org workspace variables",
			input:    []string{"myorg", "myworkspace", "variables"},
			expected: []string{"variable", "list", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace variables list",
			input:    []string{"myorg", "myworkspace", "variables", "list"},
			expected: []string{"variable", "list", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace state",
			input:    []string{"myorg", "myworkspace", "state"},
			expected: []string{"state", "list", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace state list",
			input:    []string{"myorg", "myworkspace", "state", "list"},
			expected: []string{"state", "list", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace state outputs",
			input:    []string{"myorg", "myworkspace", "state", "outputs"},
			expected: []string{"state", "outputs", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org with help flag",
			input:    []string{"myorg", "-h"},
			expected: []string{"organization:context", "-org=myorg"},
		},
		{
			name:     "org with --help flag",
			input:    []string{"myorg", "--help"},
			expected: []string{"organization:context", "-org=myorg"},
		},
		{
			name:     "org workspace with help flag",
			input:    []string{"myorg", "myworkspace", "-h"},
			expected: []string{"workspace:context", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace with --help flag",
			input:    []string{"myorg", "myworkspace", "--help"},
			expected: []string{"workspace:context", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace run-id (show)",
			input:    []string{"myorg", "myworkspace", "run-abc123"},
			expected: []string{"run", "show", "-id=run-abc123"},
		},
		{
			name:     "org workspace run-id apply",
			input:    []string{"myorg", "myworkspace", "run-abc123", "apply"},
			expected: []string{"run", "apply", "-id=run-abc123"},
		},
		{
			name:     "org workspace resources",
			input:    []string{"myorg", "myworkspace", "resources"},
			expected: []string{"workspaceresource", "list", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace tags",
			input:    []string{"myorg", "myworkspace", "tags"},
			expected: []string{"workspacetag", "list", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace configversions",
			input:    []string{"myorg", "myworkspace", "configversions"},
			expected: []string{"configversion", "list", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace assessments",
			input:    []string{"myorg", "myworkspace", "assessments"},
			expected: []string{"assessmentresult", "list", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace changerequests",
			input:    []string{"myorg", "myworkspace", "changerequests"},
			expected: []string{"changerequest", "list", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace run-id plan (shorter syntax)",
			input:    []string{"myorg", "myworkspace", "run-abc123", "plan"},
			expected: []string{"plan", "read", "-id=run-abc123"},
		},
		{
			name:     "org workspace run-id plan help (shorter syntax)",
			input:    []string{"myorg", "myworkspace", "run-abc123", "plan", "-h"},
			expected: []string{"plan", "read", "-id=run-abc123", "-h"},
		},
		{
			name:     "org workspace run-id logs (shorter syntax)",
			input:    []string{"myorg", "myworkspace", "run-abc123", "logs"},
			expected: []string{"plan", "logs", "-id=run-abc123"},
		},
		{
			name:     "org workspace run-id comments (shorter syntax)",
			input:    []string{"myorg", "myworkspace", "run-abc123", "comments"},
			expected: []string{"comment", "list", "-run-id=run-abc123"},
		},
		{
			name:     "org workspace run-id policychecks (shorter syntax)",
			input:    []string{"myorg", "myworkspace", "run-abc123", "policychecks"},
			expected: []string{"policycheck", "list", "-run-id=run-abc123"},
		},
		{
			name:     "org workspace run-id outputs (shorter syntax)",
			input:    []string{"myorg", "myworkspace", "run-abc123", "outputs"},
			expected: []string{"state", "outputs", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace run-id state (shorter syntax)",
			input:    []string{"myorg", "myworkspace", "run-abc123", "state"},
			expected: []string{"state", "list", "-org=myorg", "-workspace=myworkspace"},
		},
		{
			name:     "org workspace run-id configversion (shorter syntax)",
			input:    []string{"myorg", "myworkspace", "run-abc123", "configversion"},
			expected: []string{"configversion", "read", "-run-id=run-abc123"},
		},
		{
			name:     "org workspace runs run-id comments (longer syntax)",
			input:    []string{"myorg", "myworkspace", "runs", "run-abc123", "comments"},
			expected: []string{"comment", "list", "-run-id=run-abc123"},
		},
		{
			name:     "org workspace runs run-id comments help (longer syntax)",
			input:    []string{"myorg", "myworkspace", "runs", "run-abc123", "comments", "-h"},
			expected: []string{"comment", "list", "-run-id=run-abc123", "-h"},
		},
		{
			name:     "org workspace runs run-id policychecks (longer syntax)",
			input:    []string{"myorg", "myworkspace", "runs", "run-abc123", "policychecks"},
			expected: []string{"policycheck", "list", "-run-id=run-abc123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := r.TranslateArgs(tt.input)
			if err != nil {
				t.Errorf("TranslateArgs() error = %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("TranslateArgs() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsKnownCommand(t *testing.T) {
	r := newTestRouter()

	tests := []struct {
		arg      string
		expected bool
	}{
		{"workspace", true},
		{"run", true},
		{"organization", true},
		{"login", true},
		{"logout", true},
		{"version", true},
		{"whoami", true},
		{"notacommand", false},
		{"myorg", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			result := r.isKnownCommand(tt.arg)
			if result != tt.expected {
				t.Errorf("isKnownCommand(%q) = %v, want %v", tt.arg, result, tt.expected)
			}
		})
	}
}

func TestHasHelpFlag(t *testing.T) {
	r := newTestRouter()

	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{"no flags", []string{"myorg", "myworkspace"}, false},
		{"short help flag", []string{"myorg", "-h"}, true},
		{"long help flag", []string{"myorg", "--help"}, true},
		{"help flag", []string{"myorg", "-help"}, true},
		{"help in middle", []string{"myorg", "-h", "workspace"}, true},
		{"help at end", []string{"myorg", "workspace", "--help"}, true},
		{"empty args", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.hasHelpFlag(tt.args)
			if result != tt.expected {
				t.Errorf("hasHelpFlag(%v) = %v, want %v", tt.args, result, tt.expected)
			}
		})
	}
}

func TestIsResourceKeyword(t *testing.T) {
	r := newTestRouter()

	tests := []struct {
		arg      string
		expected bool
	}{
		{"workspaces", true},
		{"projects", true},
		{"teams", true},
		{"policies", true},
		{"policysets", true},
		{"runs", true},
		{"variables", true},
		{"state", true},
		{"workspace", false},
		{"myworkspace", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			result := r.isResourceKeyword(tt.arg)
			if result != tt.expected {
				t.Errorf("isResourceKeyword(%q) = %v, want %v", tt.arg, result, tt.expected)
			}
		})
	}
}

func TestIsKnownCommandWithInjectedCommands(t *testing.T) {
	r := NewRouter(nil, []string{"alpha", "beta echo"})

	if !r.isKnownCommand("alpha") {
		t.Fatal("expected injected command alpha to be known")
	}
	if !r.isKnownCommand("beta") {
		t.Fatal("expected injected command beta to be known")
	}
	if r.isKnownCommand("gamma") {
		t.Fatal("expected gamma to be unknown")
	}
}
