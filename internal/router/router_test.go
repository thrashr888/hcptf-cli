package router

import (
	"reflect"
	"testing"
)

func TestTranslateArgs(t *testing.T) {
	r := NewRouter(nil)

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
			name:     "org only",
			input:    []string{"myorg"},
			expected: []string{"organization", "show", "-org=myorg"},
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
			expected: []string{"workspace", "read", "-org=myorg", "-workspace=myworkspace"},
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
	r := NewRouter(nil)

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
